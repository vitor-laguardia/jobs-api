package job

import (
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/user"
)

func TestMemoryRepository(t *testing.T) {
	RepositoryContract{
		NewRepositories: func() (Repository, user.Repository) {
			userRepo := user.NewInMemoryRepository()
			return NewInMemoryRepository(userRepo), userRepo
		},
	}.Test(t)
}
