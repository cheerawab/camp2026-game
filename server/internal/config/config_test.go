package config

import "testing"

func TestLoadContentDirDefault(t *testing.T) {
	t.Setenv("CONTENT_DIR", "")
	t.Setenv("LOG_LEVEL", "info")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.ContentDir != defaultContentDir {
		t.Fatalf("expected default content dir %q, got %q", defaultContentDir, cfg.ContentDir)
	}
}

func TestLoadContentDirOverride(t *testing.T) {
	t.Setenv("CONTENT_DIR", "/tmp/camp2026-content")
	t.Setenv("LOG_LEVEL", "info")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.ContentDir != "/tmp/camp2026-content" {
		t.Fatalf("expected content dir override, got %q", cfg.ContentDir)
	}
}
