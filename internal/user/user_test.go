package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vitor-laguardia/jobs-api/internal/shared/api"
	"github.com/vitor-laguardia/jobs-api/internal/shared/assert"
)

type StubRepository struct {
	users map[string]User
}

// TODO: change to sr StubRepository
func (sr *StubRepository) GetByID(userID string) (User, error) {

	user, exists := sr.users[userID]

	if !exists {
		return User{}, ErrNotFound
	}
	return user, nil
}

func (sr *StubRepository) Create(user User) User {
	// TODO check nil
	sr.users[user.ID] = user
	return user
}

func TestGET(t *testing.T) {
	testUserID := "ab342-sbfdau-adufba-audbfuda"
	repo := newStubRepository(testUserID)
	service := newService(repo)
	handler := NewHandler(service)

	t.Run("get expected user", func(t *testing.T) {
		req := newGETUserHTTPRequest(testUserID)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got User
		expected := repo.users[testUserID]

		assert.Equal(t, res.Code, http.StatusOK, "did not get correct response status code")
		assert.Equal(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
		assert.JSONDecode(t, res.Body, &got)
		assert.Equal(t, got.Name, expected.Name, "did not get correct user Name")
		assert.Equal(t, got.ID, expected.ID, "did not get correct userID")
		assert.Equal(t, got.Email, expected.Email, "did not get correct user Email")
		assertTimeEqual(t, got.CreatedAt, expected.CreatedAt, "did not get correct user CreatedAt")
	})

	t.Run("return error when user not found", func(t *testing.T) {
		fakeUserID := "xa142-sbfdau-adufba-audbfuda"
		req := newGETUserHTTPRequest(fakeUserID)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got api.ErrorResponse

		assert.Equal(t, res.Code, http.StatusNotFound, "did not get correct response status code")
		assert.Equal(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
		assert.JSONDecode(t, res.Body, &got)
		assert.Equal(t, got.Message, MsgNotFound, "did not get correct Message in ErrorResponse")
		assert.Equal(t, got.Message, MsgNotFound, "did not get correct Message in ErrorResponse")
	})
}

func TestPOST(t *testing.T) {
	//testUserID := "cb322-sefdau-adufba-audbfuda"

	t.Run("successfuly POST User", func(t *testing.T) {
		repo := &StubRepository{make(map[string]User)}
		service := newService(repo)
		handler := NewHandler(service)

		rawPayload := `{"name":"claw", "email":"claw@gmail.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(rawPayload))
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got User
		var expected CreateUserRequest

		assert.Equal(t, res.Code, http.StatusCreated, "did not get correct response status code")
		assert.JSONDecode(t, res.Body, &got)
		assert.Unmarshal(t, rawPayload, &expected)
		assert.Equal(t, got.Name, expected.Name, "did not get correct Name in response body")
		assert.Equal(t, got.Email, expected.Email, "did not get correct Email in response body")

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
		users: map[string]User{
			userID: User{
				ID:        userID,
				Name:      "alfred",
				Email:     "alfred@gmail.com",
				CreatedAt: time.Now(),
			},
		},
	}
}

func assertTimeEqual(t *testing.T, got, want time.Time, context string) {
	t.Helper()

	if !got.Equal(want) {
		t.Errorf("%s, got: %v, want %v", context, got, want)
	}
}
