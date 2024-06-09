package service

import (
	"context"
	"fmt"

	"github.com/mickamy/go_todo_app/auth"
	"github.com/mickamy/go_todo_app/entity"
	"github.com/mickamy/go_todo_app/store"
)

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (t *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("cannot get user id from ctx")
	}
	ts, err := t.Repo.ListTasks(ctx, t.DB, id)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return ts, nil
}
