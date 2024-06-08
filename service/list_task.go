package service

import (
	"context"
	"fmt"

	"github.com/mickamy/go_todo_app/entity"
	"github.com/mickamy/go_todo_app/store"
)

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (t *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	ts, err := t.Repo.ListTasks(ctx, t.DB)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return ts, nil
}
