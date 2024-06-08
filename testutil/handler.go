package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()

	var jw, jg any
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("json.Unmarshal(want) %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jg); err != nil {
		t.Fatalf("json.Unmarshal(got) %q: %v", got, err)
	}

	if diff := cmp.Diff(jw, jg); diff != "" {
		t.Errorf("got different JSON:\n%s", diff)
	}
}

func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()
	t.Cleanup(func() { _ = got.Body.Close() })
	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatalf("io.ReadAll(%q): %v", got.Request.URL, err)
	}
	if got.StatusCode != status {
		t.Fatalf("got status %d, want %d: %q", got.StatusCode, status, gb)
	}

	if len(gb) == 0 && len(body) == 0 {
		return
	}
	AssertJSON(t, body, gb)
}

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%q): %v", path, err)
	}
	return bt
}
