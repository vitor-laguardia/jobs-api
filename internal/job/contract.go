package job

import (
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/shared/assert"
	"github.com/vitor-laguardia/jobs-api/internal/user"
)

type RepositoryContract struct {
	NewRepositories func() (Repository, user.Repository)
}

func (rc RepositoryContract) Test(t *testing.T) {
	t.Run("Repository.Create job tests", func(t *testing.T) {
		t.Run("correctly persists a job in repository", func(t *testing.T) {
			jobRepo, userRepo := rc.NewRepositories()

			user := user.NewUser("alef", "alef@gmail.com")
			if _, err := userRepo.Create(user); err != nil {
				t.Fatalf("setup failed: could not create user in repository: %v", err)
			}

			job := NewJob("create game", "roleplay for fun", 2, user.ID)
			rowCounter, err := jobRepo.Create(&job)

			assert.Nil(t, err)
			assert.Equal(t, rowCounter, 1, "did not get correct row counter in repository.Create response")
		})
	})

	t.Run("return error when create a job with wrong userID", func(t *testing.T) {
		jobRepo, userRepo := rc.NewRepositories()

		user := user.NewUser("alef", "alef@gmail.com")
		if _, err := userRepo.Create(user); err != nil {
			t.Fatalf("setup failed: could not create user in repository: %v", err)
		}

		fakeID := "7c01f5a4-d927-4f4a-98b2-a4bd9cea868b"
		job := NewJob("create game", "", 2, fakeID)
		rowCounter, err := jobRepo.Create(&job)

		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Equal(t, rowCounter, 0, "did not get correct row counter in repository.Create response")
	})
}
