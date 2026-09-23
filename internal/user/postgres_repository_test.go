package user

import (
	"context"
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/shared/testutil"
)

func TestPostgreRepository(t *testing.T) {
	ctx := context.Background()
	pgContainer := testutil.CreatePostgresContainer(t, ctx)
	db := testutil.SetupDB(t, pgContainer.ConnectionString)

	RepositoryContract{
		NewRepository: func() Repository {
			testutil.CleanDB(t, db)
			return NewPostgresRepository(db)
		},
	}.Test(t)
}
