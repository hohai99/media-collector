package repository

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

// NewDB opens (or creates) a SQLite database at dbPath and returns the
// connection. The database is configured with WAL journal mode and a 5-second
// busy timeout for better concurrent-read performance.
func NewDB(dbPath string) (*sql.DB, error) {
	// Ensure the parent directory exists.
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

// RunMigrations executes the embedded schema.sql to create all tables,
// then applies incremental column migrations for existing databases.
func RunMigrations(db *sql.DB) error {
	if _, err := db.Exec(schemaSQL); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	// Incremental migrations — ALTER TABLE for existing databases.
	// SQLite ignores duplicate-column errors gracefully.
	alterStmts := []string{
		`ALTER TABLE player_configs ADD COLUMN time_per_picture INTEGER NOT NULL DEFAULT 5`,
		`ALTER TABLE player_configs ADD COLUMN sound_source TEXT NOT NULL DEFAULT ''`,
	}
	for _, stmt := range alterStmts {
		_, _ = db.Exec(stmt) // ignore "duplicate column" errors
	}

	return nil
}
