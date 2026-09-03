package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/shared/assert"
)

type StubRepository struct {
	users map[string]*User
}

func (sr *StubRepository) GetByID(userID string) *User {
	user := sr.users[userID]
	userCopy := *user
	return &userCopy
}

func TestGETUser(t *testing.T) {
	testUserID := "ab342-sbfdau-adufba-audbfuda"
	t.Run("get expected user", func(t *testing.T) {

		repo := newStubRepository(testUserID)
		service := newService(repo)
		handler := NewHandler(service)

		req := newGETUserHTTPRequest(testUserID)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got User
		expected := *repo.users[testUserID]

		assert.Equal(t, res.Code, http.StatusOK, "did not get correct response status code")
		assert.JSONDecode(t, res.Body, &got)
		assert.Equal(t, got.Name, expected.Name, "did not get correct user Name")
		assert.Equal(t, got.ID, expected.ID, "did not get correct userID")
		assert.Equal(t, got.Email, expected.Email, "did not get correct user Email")
	})
}

func newGETUserHTTPRequest(userID string) *http.Request {
	path := fmt.Sprintf("/users/%s", userID)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue("id", userID)
	return req
}

func newStubRepository(userID string) *StubRepository {
	return &StubRepository{
		users: map[string]*User{
			userID: &User{
				ID:    userID,
				Name:  "alfred",
				Email: "alfred@gmail.com",
			},
		},
	}
}
