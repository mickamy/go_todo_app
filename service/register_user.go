package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/mickamy/go_todo_app/entity"
	"github.com/mickamy/go_todo_app/store"
)

type RegisterUser struct {
	DB   store.Executor
	Repo UserRegisterer
}

func (r *RegisterUser) RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	u := entity.User{Name: name, Password: string(pw), Role: role}
	err = r.Repo.RegisterUser(ctx, r.DB, &u)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return &u, nil
}
