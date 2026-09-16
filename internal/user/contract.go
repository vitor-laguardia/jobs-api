package user

import (
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/shared/assert"
)

type RepositoryContract struct {
	NewRepository func() Repository
}

func (rc RepositoryContract) Test(t *testing.T) {
	t.Run("Repository.Create user tests", func(t *testing.T) {
		t.Run("persists new user in repository", func(t *testing.T) {
			repo := rc.NewRepository()
			expected := NewUser("charles", "charles@gmail.com")

			got, err := repo.Create(expected)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, got.ID, expected.ID, "did not get correct user ID in repository response")
			assert.Equal(t, got.Name, expected.Name, "did not get correct user Name in repository response")
			assert.Equal(t, got.Email, expected.Email, "did not get correct user Email in repository response")
			assertTimeEqual(t, got.CreatedAt, expected.CreatedAt, "did not get correct user CreatedAt in repository response")
		})

		t.Run("return error if create user with duplicate email", func(t *testing.T) {
			repo := rc.NewRepository()
			expected := NewUser("charles", "charles@gmail.com")

			repo.Create(expected)
			_, err := repo.Create(expected)
			assert.ErrorIs(t, err, ErrDuplicateEmail)
		})
	})
}
