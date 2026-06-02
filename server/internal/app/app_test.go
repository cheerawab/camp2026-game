package app

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFailsWhenContentCannotLoad(t *testing.T) {
	t.Setenv("CONTENT_DIR", filepath.Join(t.TempDir(), "missing"))

	_, err := New(context.Background())
	if err == nil {
		t.Fatal("expected content load error")
	}
	if !strings.Contains(err.Error(), "load content") {
		t.Fatalf("expected content load error, got %v", err)
	}
}
