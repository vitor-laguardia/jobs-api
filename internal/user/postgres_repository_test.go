package user

import (
	"testing"

	"github.com/vitor-laguardia/jobs-api/database"
)

func TestPostgreRepository(t *testing.T) {
	dsn := "postgres://postgres:gojobs@localhost/jobs"
	db, _ := database.New(dsn)
	t.Cleanup(func() { db.Close() })
	RepositoryContract{
		NewRepository: func() Repository {
			return NewPostgresRepository(db)
		},
	}.Test(t)
}
