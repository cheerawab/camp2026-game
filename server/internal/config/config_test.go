package config

import (
	"strings"
	"testing"
)

func TestLoadContentDirDefault(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	t.Setenv("CONTENT_DIR", "")
	t.Setenv("ADMIN_COOKIE_SECURE", "")
	t.Setenv("MONGODB_URI", "")
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
	t.Setenv("APP_ENV", "local")
	t.Setenv("CONTENT_DIR", "/tmp/camp2026-content")
	t.Setenv("ADMIN_COOKIE_SECURE", "")
	t.Setenv("MONGODB_URI", "")
	t.Setenv("LOG_LEVEL", "info")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.ContentDir != "/tmp/camp2026-content" {
		t.Fatalf("expected content dir override, got %q", cfg.ContentDir)
	}
}

func TestLoadAdminCookieSecureDefaultsFalseForLocal(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	t.Setenv("ADMIN_COOKIE_SECURE", "")
	t.Setenv("MONGODB_URI", "")
	t.Setenv("LOG_LEVEL", "info")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.AdminCookieSecure {
		t.Fatalf("expected admin cookie secure default to be false for local")
	}
}

func TestLoadAdminCookieSecureDefaultsTrueOutsideLocal(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("ADMIN_COOKIE_SECURE", "")
	t.Setenv("MONGODB_URI", "mongodb://db.example:27017/camp2026")
	t.Setenv("LOG_LEVEL", "info")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if !cfg.AdminCookieSecure {
		t.Fatalf("expected admin cookie secure default to be true for production")
	}
}

func TestLoadAdminCookieSecureOverride(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("ADMIN_COOKIE_SECURE", "false")
	t.Setenv("MONGODB_URI", "mongodb://db.example:27017/camp2026")
	t.Setenv("LOG_LEVEL", "info")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.AdminCookieSecure {
		t.Fatalf("expected ADMIN_COOKIE_SECURE=false to override production default")
	}
}

func TestLoadAdminCookieSecureRejectsInvalidOverride(t *testing.T) {
	t.Setenv("ADMIN_COOKIE_SECURE", "sometimes")
	t.Setenv("LOG_LEVEL", "info")

	if _, err := Load(); err == nil {
		t.Fatalf("expected invalid ADMIN_COOKIE_SECURE to fail config load")
	}
}

func TestLoadMongoURIDefaultsForLocal(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	t.Setenv("ADMIN_COOKIE_SECURE", "")
	t.Setenv("MONGODB_URI", "")
	t.Setenv("LOG_LEVEL", "info")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.MongoURI != defaultMongoURI {
		t.Fatalf("expected default mongo uri %q, got %q", defaultMongoURI, cfg.MongoURI)
	}
	if strings.Contains(cfg.MongoURI, "@") {
		t.Fatalf("expected default mongo uri to omit credentials, got %q", cfg.MongoURI)
	}
}

func TestLoadRequiresMongoURIForProductionEnv(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("ADMIN_COOKIE_SECURE", "")
	t.Setenv("MONGODB_URI", "")
	t.Setenv("LOG_LEVEL", "info")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected missing MONGODB_URI to fail config load")
	}
	if !strings.Contains(err.Error(), "MONGODB_URI is required") {
		t.Fatalf("expected MONGODB_URI validation error, got %v", err)
	}
}
