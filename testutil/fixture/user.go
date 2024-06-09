package fixture

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/mickamy/go_todo_app/entity"
)

func User(u *entity.User) *entity.User {
	result := &entity.User{
		ID:         entity.UserID(rand.Int()),
		Name:       "mickamy" + strconv.Itoa(rand.Int())[:5],
		Password:   "password",
		Role:       "admin",
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	if u == nil {
		return result
	}
	if u.ID != 0 {
		result.ID = u.ID
	}
	if u.Name != "" {
		result.Name = u.Name
	}
	if u.Password != "" {
		result.Password = u.Password
	}
	if u.Role != "" {
		result.Role = u.Role
	}
	if !u.CreatedAt.IsZero() {
		result.CreatedAt = u.CreatedAt
	}
	if !u.ModifiedAt.IsZero() {
		result.ModifiedAt = u.ModifiedAt
	}
	return result
}
