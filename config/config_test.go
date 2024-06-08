package config

import (
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	wantPort := 3333
	t.Setenv("PORT", strconv.Itoa(wantPort))

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create new config: %v", err)
	}
	if got.Port != wantPort {
		t.Errorf("got port %d, want %d", got.Port, wantPort)
	}
	wantEnv := "dev"
	if got.Env != wantEnv {
		t.Errorf("got env %s, want %s", got.Env, wantEnv)
	}
}
