package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mickamy/go_todo_app/entity"
	"github.com/mickamy/go_todo_app/testutil"
)

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)

	sut := &KVS{Cli: cli}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_ok"
		uid := entity.UserID(1234)
		ctx := context.Background()
		cli.Set(ctx, key, int64(uid), 30*time.Minute)
		t.Cleanup(func() {
			cli.Del(ctx, key)
		})
		got, err := sut.Load(ctx, key)
		if err != nil {
			t.Fatalf("failed to load key %s: %v", key, err)
		}
		if got != uid {
			t.Errorf("got %d, want %d", got, uid)
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Save_not_found"
		ctx := context.Background()
		got, err := sut.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound(value = %d", err, got)
		}
	})
}
