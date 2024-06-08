package config

import (
	"testing"
)

func TestNew(t *testing.T) {
	wantPort := "3333"
	t.Setenv("PORT", wantPort)

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create new config: %v", err)
	}
	if got.Port != wantPort {
		t.Errorf("got port %s, want %s", got.Port, wantPort)
	}
	wantEnv := "dev"
	if got.Env != wantEnv {
		t.Errorf("got env %s, want %s", got.Env, wantEnv)
	}
}
