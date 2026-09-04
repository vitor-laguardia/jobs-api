package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vitor-laguardia/jobs-api/internal/shared/assert"
	"github.com/vitor-laguardia/jobs-api/internal/shared/httputil"
)

type StubRepository struct {
	users map[string]*User
}

func (sr *StubRepository) GetByID(userID string) (*User, error) {
	user := sr.users[userID]
	if user == nil {
		return nil, ErrNotFound
	}
	userCopy := *user
	return &userCopy, nil
}

func TestGETUser(t *testing.T) {
	testUserID := "ab342-sbfdau-adufba-audbfuda"
	repo := newStubRepository(testUserID)
	service := newService(repo)
	handler := NewHandler(service)

	t.Run("get expected user", func(t *testing.T) {
		req := newGETUserHTTPRequest(testUserID)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got User
		expected := *repo.users[testUserID]

		assert.Equal(t, res.Code, http.StatusOK, "did not get correct response status code")
		assert.Equal(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
		assert.JSONDecode(t, res.Body, &got)
		assert.Equal(t, got.Name, expected.Name, "did not get correct user Name")
		assert.Equal(t, got.ID, expected.ID, "did not get correct userID")
		assert.Equal(t, got.Email, expected.Email, "did not get correct user Email")
		timeEqual(t, got.CreatedAt, expected.CreatedAt, "did not get correct user CreatedAt")
	})

	t.Run("return error when user not found", func(t *testing.T) {
		fakeUserID := "xa142-sbfdau-adufba-audbfuda"
		req := newGETUserHTTPRequest(fakeUserID)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got httputil.ErrorResponse

		assert.Equal(t, res.Code, http.StatusNotFound, "did not get correct response status code")
		assert.Equal(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
		assert.JSONDecode(t, res.Body, &got)
		assert.Equal(t, got.Message, MsgNotFound, "did not get correct Message in ErrorResponse")
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
				ID:        userID,
				Name:      "alfred",
				Email:     "alfred@gmail.com",
				CreatedAt: time.Now(),
			},
		},
	}
}

func timeEqual(t *testing.T, got, want time.Time, context string) {
	t.Helper()

	if !got.Equal(want) {
		t.Errorf("%s, got: %v, want %v", context, got, want)
	}
}
