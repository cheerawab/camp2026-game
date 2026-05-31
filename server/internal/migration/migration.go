package migration

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const schemaMigrationsTable = "schema_migrations"

var migrationFilePattern = regexp.MustCompile(`^([0-9]+)_(.+)\.sql$`)

type Runner struct {
	Conn *pgx.Conn
	Dir  string
	Out  io.Writer
}

type File struct {
	Version int64
	Name    string
	Path    string
	UpSQL   string
	DownSQL string
	HasDown bool
}

type AppliedMigration struct {
	Version   int64
	Name      string
	AppliedAt time.Time
}

func DefaultDir() string {
	if value := strings.TrimSpace(os.Getenv("MIGRATIONS_DIR")); value != "" {
		return value
	}
	if _, err := os.Stat("db/migrations"); err == nil {
		return "db/migrations"
	}
	return "server/db/migrations"
}

func Open(ctx context.Context, databaseURL string) (*pgx.Conn, error) {
	cfg, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	return conn, nil
}

func (r Runner) Up(ctx context.Context) error {
	if err := r.ensureReady(ctx); err != nil {
		return err
	}

	files, err := LoadFiles(r.Dir)
	if err != nil {
		return err
	}
	applied, err := r.appliedVersions(ctx)
	if err != nil {
		return err
	}

	appliedCount := 0
	for _, file := range files {
		if _, ok := applied[file.Version]; ok {
			continue
		}
		if err := r.applyUp(ctx, file); err != nil {
			return err
		}
		appliedCount++
		r.printf("applied %06d %s\n", file.Version, file.Name)
	}

	if appliedCount == 0 {
		r.printf("no pending migrations\n")
	}
	return nil
}

func (r Runner) Down(ctx context.Context) error {
	if err := r.ensureReady(ctx); err != nil {
		return err
	}

	applied, err := r.appliedList(ctx)
	if err != nil {
		return err
	}
	if len(applied) == 0 {
		r.printf("no applied migrations\n")
		return nil
	}

	files, err := LoadFiles(r.Dir)
	if err != nil {
		return err
	}
	filesByVersion := make(map[int64]File, len(files))
	for _, file := range files {
		filesByVersion[file.Version] = file
	}

	last := applied[len(applied)-1]
	file, ok := filesByVersion[last.Version]
	if !ok {
		return fmt.Errorf("migration file for applied version %06d is missing", last.Version)
	}
	if !file.HasDown || strings.TrimSpace(file.DownSQL) == "" {
		return fmt.Errorf("migration %06d %s has no down section", file.Version, file.Name)
	}

	if err := r.applyDown(ctx, file); err != nil {
		return err
	}
	r.printf("reverted %06d %s\n", file.Version, file.Name)
	return nil
}

func (r Runner) Status(ctx context.Context) error {
	if err := r.ensureReady(ctx); err != nil {
		return err
	}

	files, err := LoadFiles(r.Dir)
	if err != nil {
		return err
	}
	applied, err := r.appliedVersions(ctx)
	if err != nil {
		return err
	}

	seen := make(map[int64]struct{}, len(files))
	for _, file := range files {
		seen[file.Version] = struct{}{}
		state := "pending"
		if _, ok := applied[file.Version]; ok {
			state = "applied"
		}
		r.printf("%-8s %06d %s\n", state, file.Version, file.Name)
	}

	var missing []AppliedMigration
	for _, migration := range applied {
		if _, ok := seen[migration.Version]; !ok {
			missing = append(missing, migration)
		}
	}
	sort.Slice(missing, func(i, j int) bool {
		return missing[i].Version < missing[j].Version
	})
	for _, migration := range missing {
		r.printf("%-8s %06d %s\n", "missing", migration.Version, migration.Name)
	}
	return nil
}

func LoadFiles(dir string) ([]File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}

	var files []File
	versions := map[int64]string{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		matches := migrationFilePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}

		version, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse migration version %q: %w", entry.Name(), err)
		}
		if previous, ok := versions[version]; ok {
			return nil, fmt.Errorf("duplicate migration version %06d: %s and %s", version, previous, entry.Name())
		}
		versions[version] = entry.Name()

		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}

		upSQL, downSQL, hasDown, err := ParseSQL(string(content))
		if err != nil {
			return nil, fmt.Errorf("parse migration %s: %w", entry.Name(), err)
		}
		if strings.TrimSpace(upSQL) == "" {
			return nil, fmt.Errorf("migration %s has empty up section", entry.Name())
		}

		files = append(files, File{
			Version: version,
			Name:    strings.TrimSuffix(matches[2], ".sql"),
			Path:    path,
			UpSQL:   upSQL,
			DownSQL: downSQL,
			HasDown: hasDown,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Version < files[j].Version
	})
	return files, nil
}

func ParseSQL(content string) (upSQL string, downSQL string, hasDown bool, err error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	section := "up"
	seenUpMarker := false
	var upLines []string
	var downLines []string

	for scanner.Scan() {
		line := scanner.Text()
		marker := strings.ToLower(strings.TrimSpace(line))
		switch marker {
		case "-- +goose up":
			if seenUpMarker {
				return "", "", false, errors.New("duplicate up marker")
			}
			seenUpMarker = true
			section = "up"
			continue
		case "-- +goose down":
			if hasDown {
				return "", "", false, errors.New("duplicate down marker")
			}
			hasDown = true
			section = "down"
			continue
		}
		if strings.HasPrefix(marker, "-- +goose") {
			continue
		}

		switch section {
		case "up":
			upLines = append(upLines, line)
		case "down":
			downLines = append(downLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", false, err
	}

	return strings.TrimSpace(strings.Join(upLines, "\n")), strings.TrimSpace(strings.Join(downLines, "\n")), hasDown, nil
}

func CreateFile(dir, name string) (string, error) {
	slug, err := slugName(name)
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read migrations dir: %w", err)
	}

	nextVersion := int64(1)
	versions := map[int64]string{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		matches := migrationFilePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}
		version, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return "", fmt.Errorf("parse migration version %q: %w", entry.Name(), err)
		}
		if previous, ok := versions[version]; ok {
			return "", fmt.Errorf("duplicate migration version %06d: %s and %s", version, previous, entry.Name())
		}
		versions[version] = entry.Name()
		if version >= nextVersion {
			nextVersion = version + 1
		}
	}

	path := filepath.Join(dir, fmt.Sprintf("%06d_%s.sql", nextVersion, slug))
	template := "-- +goose Up\n\n-- +goose Down\n"
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", fmt.Errorf("write migration file: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString(template); err != nil {
		return "", fmt.Errorf("write migration file: %w", err)
	}
	return path, nil
}

func (r Runner) ensureReady(ctx context.Context) error {
	if r.Conn == nil {
		return errors.New("migration connection is nil")
	}
	if strings.TrimSpace(r.Dir) == "" {
		return errors.New("migrations directory is required")
	}
	if _, err := r.Conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version bigint PRIMARY KEY,
    name text NOT NULL,
    applied_at timestamptz NOT NULL DEFAULT now()
)`); err != nil {
		return fmt.Errorf("ensure %s table: %w", schemaMigrationsTable, err)
	}
	return nil
}

func (r Runner) appliedVersions(ctx context.Context) (map[int64]AppliedMigration, error) {
	applied, err := r.appliedList(ctx)
	if err != nil {
		return nil, err
	}
	byVersion := make(map[int64]AppliedMigration, len(applied))
	for _, migration := range applied {
		byVersion[migration.Version] = migration
	}
	return byVersion, nil
}

func (r Runner) appliedList(ctx context.Context) ([]AppliedMigration, error) {
	rows, err := r.Conn.Query(ctx, `
SELECT version, name, applied_at
FROM schema_migrations
ORDER BY version ASC`)
	if err != nil {
		return nil, fmt.Errorf("list applied migrations: %w", err)
	}
	defer rows.Close()

	var applied []AppliedMigration
	for rows.Next() {
		var migration AppliedMigration
		if err := rows.Scan(&migration.Version, &migration.Name, &migration.AppliedAt); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}
		applied = append(applied, migration)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}
	return applied, nil
}

func (r Runner) applyUp(ctx context.Context, file File) error {
	tx, err := r.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration %06d: %w", file.Version, err)
	}
	defer rollback(ctx, tx)

	if _, err := tx.Exec(ctx, file.UpSQL); err != nil {
		return fmt.Errorf("apply migration %06d %s: %w", file.Version, file.Name, err)
	}
	if _, err := tx.Exec(ctx, `
INSERT INTO schema_migrations (version, name)
VALUES ($1, $2)`, file.Version, file.Name); err != nil {
		return fmt.Errorf("record migration %06d %s: %w", file.Version, file.Name, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %06d %s: %w", file.Version, file.Name, err)
	}
	return nil
}

func (r Runner) applyDown(ctx context.Context, file File) error {
	tx, err := r.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin rollback %06d: %w", file.Version, err)
	}
	defer rollback(ctx, tx)

	if _, err := tx.Exec(ctx, file.DownSQL); err != nil {
		return fmt.Errorf("revert migration %06d %s: %w", file.Version, file.Name, err)
	}
	if _, err := tx.Exec(ctx, `
DELETE FROM schema_migrations
WHERE version = $1`, file.Version); err != nil {
		return fmt.Errorf("remove migration record %06d %s: %w", file.Version, file.Name, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit rollback %06d %s: %w", file.Version, file.Name, err)
	}
	return nil
}

func (r Runner) printf(format string, args ...any) {
	if r.Out == nil {
		return
	}
	_, _ = fmt.Fprintf(r.Out, format, args...)
}

func rollback(ctx context.Context, tx pgx.Tx) {
	_ = tx.Rollback(ctx)
}

func slugName(name string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(name))
	if value == "" {
		return "", errors.New("migration name is required")
	}

	var builder strings.Builder
	lastUnderscore := false
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			builder.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore {
			builder.WriteByte('_')
			lastUnderscore = true
		}
	}

	slug := strings.Trim(builder.String(), "_")
	if slug == "" {
		return "", errors.New("migration name must contain letters or numbers")
	}
	return slug, nil
}
