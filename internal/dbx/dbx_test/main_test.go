package dbx

import (
	"database/sql"
	"testing"
)

const databaseURL = "postgresql://postgres:postgres@localhost:7777/postgres?sslmode=disable"

func strPtr(s string) *string { return &s }

func openDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}
	return db
}

func mustExec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatalf("exec failed: %v", err)
	}
}
