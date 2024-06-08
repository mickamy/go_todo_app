package service

import (
	"context"

	"github.com/mickamy/go_todo_app/entity"
	"github.com/mickamy/go_todo_app/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskLister UserRegisterer
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Executor, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}

type UserRegisterer interface {
	RegisterUser(ctx context.Context, db store.Executor, u *entity.User) error
}
