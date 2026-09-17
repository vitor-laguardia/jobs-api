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
				t.Fatalf("setup failed: could not create user in repository: %v", err)
			}
			assert.Equal(t, got.ID, expected.ID, "did not get correct user ID in repository response")
			assert.Equal(t, got.Name, expected.Name, "did not get correct user Name in repository response")
			assert.Equal(t, got.Email, expected.Email, "did not get correct user Email in repository response")
			assert.TimeEqual(t, got.CreatedAt, expected.CreatedAt, "did not get correct user CreatedAt in repository response")
		})

		t.Run("return error if create user with duplicate email", func(t *testing.T) {
			repo := rc.NewRepository()
			expected := NewUser("charles", "charles@gmail.com")

			if _, err := repo.Create(expected); err != nil {
				t.Fatalf("setup failed: could not create user in repository: %v", err)
			}
			_, err := repo.Create(expected)

			assert.ErrorIs(t, err, ErrDuplicateEmail)
		})
	})
	t.Run("repository.GetByID user tests", func(t *testing.T) {
		t.Run("sucessfuly get user", func(t *testing.T) {
			repo := rc.NewRepository()
			expected := NewUser("charles", "charles@gmail.com")

			if _, err := repo.Create(expected); err != nil {
				t.Fatalf("setup failed: could not create user in repository: %v", err)
			}

			got, err := repo.GetByID(expected.ID)
			if err != nil {
				t.Fatalf("setup failed: could not getByID in repository: %v", err)
			}

			assert.Equal(t, got.ID, expected.ID, "did not get correct user ID in repository response")
			assert.Equal(t, got.Name, expected.Name, "did not get correct user Name in repository response")
			assert.Equal(t, got.Email, expected.Email, "did not get correct user Email in repository response")
			assert.TimeEqual(t, got.CreatedAt, expected.CreatedAt, "did not get correct user CreatedAt in repository response")
		})

		t.Run("return error when user not found", func(t *testing.T) {
			repo := rc.NewRepository()
			expected := NewUser("other user", "other_user@gmail.com")

			_, err := repo.GetByID(expected.ID)

			assert.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("repository.Update user tests", func(t *testing.T) {
		t.Run("sucessfuly update user", func(t *testing.T) {
			repo := rc.NewRepository()
			expected := NewUser("alex", "alex@gmail.com")

			if _, err := repo.Create(expected); err != nil {
				t.Fatalf("setup failed: could not create user in repository: %v", err)
			}

			expected.Name = "gon"

			got, err := repo.Update(expected)
			if err != nil {
				t.Fatalf("setup failed: could not update user in repository: %v", err)
			}

			assert.Equal(t, got.ID, expected.ID, "did not get correct user ID in repository response")
			assert.Equal(t, got.Name, expected.Name, "did not get correct user Name in repository response")
			assert.Equal(t, got.Email, expected.Email, "did not get correct user Email in repository response")
			assert.TimeEqual(t, got.CreatedAt, expected.CreatedAt, "did not get correct user CreatedAt in repository response")
		})
	})

	t.Run("repository.Delete user tests", func(t *testing.T) {
		t.Run("sucessfuly delete user", func(t *testing.T) {
			repo := rc.NewRepository()
			expected := NewUser("alex", "alex@gmail.com")

			if _, err := repo.Create(expected); err != nil {
				t.Fatalf("setup failed: could not create user in repository: %v", err)
			}

			err := repo.Delete(expected.ID)
			if err != nil {
				t.Fatalf("setup failed: could not delete user in repository: %v", err)
			}

			got, err := repo.GetByID(expected.ID)

			assert.ErrorIs(t, err, ErrNotFound)
			assert.Equal(t, got, User{}, "did not get correct user Email in repository response")
		})
	})
	t.Run("return error when user not found", func(t *testing.T) {
		repo := rc.NewRepository()
		fakeID := "264664e5-d784-476b-adfb-341cd4e64c2b"

		err := repo.Delete(fakeID)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}
