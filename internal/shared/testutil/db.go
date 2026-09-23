package testutil

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func SetupDB(t *testing.T, dsn string) *sql.DB {
	t.Helper()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("error sql.Open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("error db.Ping: %v", err)
	}
	return db
}

func CleanDB(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec("TRUNCATE TABLE jobs, users CASCADE")
	if err != nil {
		t.Fatalf("error TRUNCATE TABLE: %v", err)
	}
}
