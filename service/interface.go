package service

import (
	"context"

	"github.com/mickamy/go_todo_app/entity"
	"github.com/mickamy/go_todo_app/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskLister UserRegisterer UserGetter TokenGenerator
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Executor, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer, userID entity.UserID) (entity.Tasks, error)
}

type UserRegisterer interface {
	RegisterUser(ctx context.Context, db store.Executor, u *entity.User) error
}

type UserGetter interface {
	GetUser(ctx context.Context, db store.Queryer, id string) (*entity.User, error)
}

type TokenGenerator interface {
	GenerateToken(ctx context.Context, user entity.User) ([]byte, error)
}
