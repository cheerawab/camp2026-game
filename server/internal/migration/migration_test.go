package migration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSQLWithGooseSections(t *testing.T) {
	content := `-- +goose Up
CREATE TABLE players (id uuid PRIMARY KEY);
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd

-- +goose Down
DROP TABLE players;
`

	upSQL, downSQL, hasDown, err := ParseSQL(content)
	if err != nil {
		t.Fatalf("ParseSQL returned error: %v", err)
	}
	if !hasDown {
		t.Fatal("expected down section")
	}
	if upSQL != "CREATE TABLE players (id uuid PRIMARY KEY);\nSELECT 1;" {
		t.Fatalf("unexpected up SQL:\n%s", upSQL)
	}
	if downSQL != "DROP TABLE players;" {
		t.Fatalf("unexpected down SQL:\n%s", downSQL)
	}
}

func TestParseSQLWithoutMarkersUsesWholeFileAsUp(t *testing.T) {
	upSQL, downSQL, hasDown, err := ParseSQL("SELECT 1;\n")
	if err != nil {
		t.Fatalf("ParseSQL returned error: %v", err)
	}
	if upSQL != "SELECT 1;" {
		t.Fatalf("unexpected up SQL: %q", upSQL)
	}
	if downSQL != "" {
		t.Fatalf("unexpected down SQL: %q", downSQL)
	}
	if hasDown {
		t.Fatal("expected no down section")
	}
}

func TestLoadFilesSortsAndRejectsDuplicateVersions(t *testing.T) {
	dir := t.TempDir()
	writeTestMigration(t, dir, "000002_second.sql")
	writeTestMigration(t, dir, "000001_first.sql")

	files, err := LoadFiles(dir)
	if err != nil {
		t.Fatalf("LoadFiles returned error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[0].Version != 1 || files[0].Name != "first" {
		t.Fatalf("unexpected first migration: %#v", files[0])
	}
	if files[1].Version != 2 || files[1].Name != "second" {
		t.Fatalf("unexpected second migration: %#v", files[1])
	}

	writeTestMigration(t, dir, "000001_duplicate.sql")
	if _, err := LoadFiles(dir); err == nil {
		t.Fatal("expected duplicate version error")
	}
}

func TestCreateFileUsesNextVersionAndSlug(t *testing.T) {
	dir := t.TempDir()
	writeTestMigration(t, dir, "000001_init.sql")
	if err := os.WriteFile(filepath.Join(dir, "000002_empty.sql"), []byte("-- +goose Up\n\n-- +goose Down\n"), 0o644); err != nil {
		t.Fatalf("write empty migration: %v", err)
	}

	path, err := CreateFile(dir, "Add Players Table")
	if err != nil {
		t.Fatalf("CreateFile returned error: %v", err)
	}

	if filepath.Base(path) != "000003_add_players_table.sql" {
		t.Fatalf("unexpected path: %s", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read created migration: %v", err)
	}
	if string(content) != "-- +goose Up\n\n-- +goose Down\n" {
		t.Fatalf("unexpected file content: %q", string(content))
	}
}

func writeTestMigration(t *testing.T, dir, name string) {
	t.Helper()
	content := "-- +goose Up\nSELECT 1;\n-- +goose Down\nSELECT 2;\n"
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write test migration %s: %v", name, err)
	}
}
