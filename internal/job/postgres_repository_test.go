package job

import (
	"context"
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/shared/testutil"
	"github.com/vitor-laguardia/jobs-api/internal/user"
)

func TestPostgresRepository(t *testing.T) {
	ctx := context.Background()
	pgContainer := testutil.CreatePostgresContainer(t, ctx)
	db := testutil.SetupDB(t, pgContainer.ConnectionString)

	RepositoryContract{
		NewRepositories: func() (Repository, user.Repository) {
			testutil.CleanDB(t, db)
			return NewPostgresRepository(db), user.NewPostgresRepository(db)
		},
	}.Test(t)
}
