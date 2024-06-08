package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/mickamy/go_todo_app/entity"
	"github.com/mickamy/go_todo_app/testutil"
)

func TestAddTask(t *testing.T) {
	t.Parallel()
	type want struct {
		status  int
		resFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/created_req.golden.json",
			want: want{
				status:  http.StatusCreated,
				resFile: "testdata/add_task/created_res.golden.json",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_task/bad_req.golden.json",
			want: want{
				status:  http.StatusBadRequest,
				resFile: "testdata/add_task/bad_res.golden.json",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)
			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(ctx context.Context, title string) (*entity.Task, error) {
				if tt.want.status == http.StatusCreated {
					return &entity.Task{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}

			sut := AddTask{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			res := w.Result()
			testutil.AssertResponse(t,
				res, tt.want.status, testutil.LoadFile(t, tt.want.resFile),
			)
		})
	}
}
