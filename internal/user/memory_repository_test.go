package user

import "testing"

func TestMemoryRepository(t *testing.T) {
	RepositoryContract{
		NewRepository: func() Repository {
			return NewInMemoryRepository()
		},
	}.Test(t)
}
