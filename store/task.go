package store

import (
	"context"

	"github.com/mickamy/go_todo_app/entity"
)

func (r *Repository) ListTasks(ctx context.Context, db Queryer, userID entity.UserID) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT 
    			id, user_id, title, status, created_at, modified_at
			FROM task where user_id = ?;`
	if err := db.SelectContext(ctx, &tasks, sql, userID); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(ctx context.Context, db Executor, task *entity.Task) error {
	task.CreatedAt = r.Clocker.Now()
	task.ModifiedAt = r.Clocker.Now()
	sql := `INSERT INTO task (user_id, title, status, created_at, modified_at) VALUES (?, ?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, sql, task.UserID, task.Title, task.Status, task.CreatedAt, task.ModifiedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = entity.TaskID(id)
	return nil
}
