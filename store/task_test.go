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
)

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)
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
		`INSERT INTO task \(title, status, created_at, modified_at\) VALUES \(\?\, \?\, \?\, \?\)`,
	).WithArgs(okTask.Title, okTask.Status, okTask.CreatedAt, okTask.ModifiedAt).
		WillReturnResult(sqlmock.NewResult(int64(wantID), 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Executor) entity.Tasks {
	t.Helper()

	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("failed to delete task: %v", err)
	}
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title: "want task 1", Status: "todo",
			CreatedAt: c.Now(), ModifiedAt: c.Now(),
		},
		{
			Title: "want task 2", Status: "todo",
			CreatedAt: c.Now(), ModifiedAt: c.Now(),
		},
		{
			Title: "want task 3", Status: "done",
			CreatedAt: c.Now(), ModifiedAt: c.Now(),
		},
	}
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, status, created_at, modified_at) 
					VALUES 
					    (?, ?, ?, ?),
					    (?, ?, ?, ?),
					    (?, ?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].CreatedAt, wants[0].ModifiedAt,
		wants[1].Title, wants[1].Status, wants[1].CreatedAt, wants[1].ModifiedAt,
		wants[2].Title, wants[2].Status, wants[2].CreatedAt, wants[2].ModifiedAt,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)

	return wants
}
