package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mickamy/go_todo_app/config"
)

func TestNewMux(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	cfg, err := config.New()
	sut, cleanup, err := NewMux(context.Background(), cfg)
	if err != nil {
		t.Error("NewMux() should not return an error")
	}
	defer cleanup()
	sut.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusOK {
		t.Error("want status code 200 but got", resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	want := `{"status":"ok"}`
	if string(got) != want {
		t.Errorf("want %s, got %s", want, string(got))
	}
}
