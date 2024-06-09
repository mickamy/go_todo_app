package service

import (
	"context"
	"fmt"

	"github.com/mickamy/go_todo_app/store"
)

type Login struct {
	DB             store.Queryer
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, username string, password string) (string, error) {
	u, err := l.Repo.GetUser(ctx, l.DB, username)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if err = u.ComparePassword(password); err != nil {
		return "", fmt.Errorf("failed to compare password: %w", err)
	}
	jwt, err := l.TokenGenerator.GenerateToken(ctx, *u)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return string(jwt), nil
}
