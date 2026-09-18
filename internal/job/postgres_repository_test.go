package job

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vitor-laguardia/jobs-api/internal/user"
)

func TestPostgresRepository(t *testing.T) {
	dsn := "postgres://postgres:gojobs@localhost/jobs"
	db := setupDB(t, dsn)
	t.Cleanup(func() { db.Close() })

	RepositoryContract{
		NewRepositories: func() (Repository, user.Repository) {
			cleanDB(t, db)
			return NewPostgresRepository(db), user.NewPostgresRepository(db)
		},
	}.Test(t)
}

func setupDB(t *testing.T, dsn string) *sql.DB {
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

func cleanDB(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec("TRUNCATE TABLE jobs, users CASCADE")
	if err != nil {
		t.Fatalf("error TRUNCATE TABLE: %v", err)
	}
}
