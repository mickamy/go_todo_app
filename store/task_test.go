package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"

	"github.com/mickamy/go_todo_app/clock"
	"github.com/mickamy/go_todo_app/entity"
	"github.com/mickamy/go_todo_app/testutil"
	"github.com/mickamy/go_todo_app/testutil/fixture"
)

func prepareUser(ctx context.Context, t *testing.T, db Executor) entity.UserID {
	t.Helper()
	u := fixture.User(nil)
	sql := `INSERT INTO user (name, password, role, created_at, modified_at) VALUES (?, ?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, sql, u.Name, u.Password, u.Role, u.CreatedAt, u.ModifiedAt)
	if err != nil {
		t.Fatalf("failed to prepare user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get last insert id: %v", err)
	}
	return entity.UserID(id)
}

func prepareTasks(ctx context.Context, t *testing.T, db Executor) (entity.UserID, entity.Tasks) {
	t.Helper()
	userID := prepareUser(ctx, t, db)
	otherUserID := prepareUser(ctx, t, db)
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserID:     userID,
			Title:      "want task 1",
			Status:     entity.TaskStatusTodo,
			CreatedAt:  c.Now(),
			ModifiedAt: c.Now(),
		},
		{
			UserID:     userID,
			Title:      "want task 2",
			Status:     entity.TaskStatusDone,
			CreatedAt:  c.Now(),
			ModifiedAt: c.Now(),
		},
	}
	tasks := entity.Tasks{
		wants[0],
		{
			UserID:     otherUserID,
			Title:      "not want task",
			Status:     entity.TaskStatusTodo,
			CreatedAt:  c.Now(),
			ModifiedAt: c.Now(),
		},
		wants[1],
	}
	sql := `INSERT INTO task (user_id, title, status, created_at, modified_at) VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?), (?, ?, ?, ?, ?);`
	result, err := db.ExecContext(ctx, sql,
		tasks[0].UserID, tasks[0].Title, tasks[0].Status, tasks[0].CreatedAt, tasks[0].ModifiedAt,
		tasks[1].UserID, tasks[1].Title, tasks[1].Status, tasks[1].CreatedAt, tasks[1].ModifiedAt,
		tasks[2].UserID, tasks[2].Title, tasks[2].Status, tasks[2].CreatedAt, tasks[2].ModifiedAt,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	tasks[0].ID = entity.TaskID(id)
	tasks[1].ID = entity.TaskID(id + 1)
	tasks[2].ID = entity.TaskID(id + 2)
	return userID, wants
}

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	userID, wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx, userID)
	if err != nil {
		t.Fatalf("unexpected error listing tasks: %v", err)
	}
	if d := cmp.Diff(wants, gots); len(d) > 0 {
		t.Errorf("ListTasks mismatch (-want +got):\n%s", d)
	}
}

func TestRepository_ListTasksMock(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		Title:      "ok task",
		Status:     "todo",
		CreatedAt:  c.Now(),
		ModifiedAt: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	mock.ExpectExec(
		`INSERT INTO task \(user_id, title, status, created_at, modified_at\) VALUES \(\?\, \?\, \?\, \?\, \?\)`,
	).WithArgs(okTask.UserID, okTask.Title, okTask.Status, okTask.CreatedAt, okTask.ModifiedAt).
		WillReturnResult(sqlmock.NewResult(int64(wantID), 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}
