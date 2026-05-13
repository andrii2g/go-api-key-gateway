package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Runner struct {
	db            *sql.DB
	migrationsDir string
}

type StatusRow struct {
	Version string
	Applied bool
}

func NewRunner(db *sql.DB, migrationsDir string) *Runner {
	if migrationsDir == "" {
		migrationsDir = filepath.Join("migrations", "mysql")
	}
	return &Runner{db: db, migrationsDir: migrationsDir}
}

func (r *Runner) Up(ctx context.Context) error {
	if err := r.ensureSchemaMigrations(ctx); err != nil {
		return err
	}
	files, err := r.upFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		applied, err := r.isApplied(ctx, file)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		content, err := os.ReadFile(filepath.Join(r.migrationsDir, file))
		if err != nil {
			return err
		}
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply %s: %w", file, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES (?)`, file); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) Status(ctx context.Context) ([]StatusRow, error) {
	if err := r.ensureSchemaMigrations(ctx); err != nil {
		return nil, err
	}
	files, err := r.upFiles()
	if err != nil {
		return nil, err
	}
	rows := make([]StatusRow, 0, len(files))
	for _, file := range files {
		applied, err := r.isApplied(ctx, file)
		if err != nil {
			return nil, err
		}
		rows = append(rows, StatusRow{Version: file, Applied: applied})
	}
	return rows, nil
}

func (r *Runner) ensureSchemaMigrations(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) NOT NULL PRIMARY KEY,
    applied_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);`)
	return err
}

func (r *Runner) isApplied(ctx context.Context, version string) (bool, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM schema_migrations WHERE version = ?`, version).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Runner) upFiles() ([]string, error) {
	entries, err := os.ReadDir(r.migrationsDir)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)
	return files, nil
}
