package job

import (
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/user"
)

func TestMemoryRepository(t *testing.T) {
	RepositoryContract{
		NewRepositories: func() (Repository, user.Repository) {
			return NewInMemoryRepository(), user.NewInMemoryRepository()
		},
	}.Test(t)
}
