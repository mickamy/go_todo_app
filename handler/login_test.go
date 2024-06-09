package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/mickamy/go_todo_app/testutil"
)

func TestLogin_ServeHTTP(t *testing.T) {
	type moq struct {
		token string
		err   error
	}
	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		reqFile string
		moq     moq
		want    want
	}{
		"ok": {
			reqFile: "testdata/login/ok_req.golden.json",
			moq: moq{
				token: "from_moq",
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/login/ok_rsp.golden.json",
			},
		},
		"badRequest": {
			reqFile: "testdata/login/bad_req.golden.json",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/login/bad_rsp.golden.json",
			},
		},
		"internalServerError": {
			reqFile: "testdata/login/ok_req.golden.json",
			moq: moq{
				err: errors.New("error from mock"),
			},
			want: want{
				status:  http.StatusInternalServerError,
				rspFile: "testdata/login/internal_server_error_rsp.golden.json",
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(tt *testing.T) {
			tt.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewReader(testutil.LoadFile(t, tc.reqFile)),
			)

			moq := &LoginServiceMock{}
			moq.LoginFunc = func(ctx context.Context, email string, password string) (string, error) {
				return tc.moq.token, tc.moq.err
			}
			sut := Login{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t, resp, tc.want.status, testutil.LoadFile(t, tc.want.rspFile))
		})
	}
}
