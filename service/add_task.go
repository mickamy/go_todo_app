package service

import (
	"context"
	"fmt"

	"github.com/mickamy/go_todo_app/auth"
	"github.com/mickamy/go_todo_app/entity"
	"github.com/mickamy/go_todo_app/store"
)

type AddTask struct {
	DB   store.Executor
	Repo TaskAdder
}

func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("cannot get user id from ctx")
	}
	t := &entity.Task{
		UserID: id,
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	err := a.Repo.AddTask(ctx, a.DB, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
