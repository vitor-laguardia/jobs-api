package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vitor-laguardia/jobs-api/internal/shared/api"
	"github.com/vitor-laguardia/jobs-api/internal/shared/assert"
)

type StubRepository struct {
	users map[string]User
}

func (sr *StubRepository) GetByID(userID string) (User, error) {
	user, exists := sr.users[userID]

	if !exists {
		return User{}, ErrNotFound
	}

	return user, nil
}

func (sr *StubRepository) Create(user User) User {
	sr.users[user.ID] = user
	return user
}

func (sr *StubRepository) Update(user User) User {
	sr.users[user.ID] = user
	return user
}

func TestGET(t *testing.T) {
	testUserID := "ab342-sbfdau-adufba-audbfuda"
	repo := newStubRepository(testUserID)
	service := NewService(repo)
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
	t.Run("successfully POST User", func(t *testing.T) {
		repo := &StubRepository{users: make(map[string]User)}
		service := NewService(repo)
		handler := NewHandler(service)
		rawPayload := `{"name":"claw", "email":"claw@gmail.com"}`
		req := newPOSTUserHTTPRequest(rawPayload)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got User
		var expected CreateUserRequest

		assert.Equal(t, res.Code, http.StatusCreated, "did not get correct response status code")
		assert.Equal(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
		assert.JSONDecode(t, res.Body, &got)
		assert.Unmarshal(t, rawPayload, &expected)

		t.Run("response body is correct", func(t *testing.T) {
			assert.Equal(t, got.Name, expected.Name, "did not get correct Name in response body")
			assert.Equal(t, got.Email, expected.Email, "did not get correct Email in response body")
			assertUUIDFormat(t, got.ID)
			assertTimeFormat(t, got.CreatedAt)
		})

		t.Run("user is correctly persisted in repository", func(t *testing.T) {
			userStub := repo.users[got.ID]
			assert.Equal(t, got.Name, userStub.Name, "user Name was not correctly persisted in repository")
			assert.Equal(t, got.ID, userStub.ID, "user ID was not correctly persisted in repository")
			assert.Equal(t, got.Email, userStub.Email, "user Email was not correctly persisted in repository")
			assertTimeEqual(t, got.CreatedAt, userStub.CreatedAt, "user CreatedAt was not correctly persisted in repository")
		})
	})

	t.Run("validates disallowed payload values", func(t *testing.T) {
		repo := &StubRepository{users: make(map[string]User)}
		service := NewService(repo)
		handler := NewHandler(service)

		testCases := []struct {
			name        string
			rawPayload  string
			errorKey    string
			wantMessage string
		}{
			{
				name:        "missing name",
				rawPayload:  `{"name":"" }`,
				errorKey:    keyName,
				wantMessage: MsgNameRequired,
			},
			{
				name:        "missing email",
				rawPayload:  `{"email":"" }`,
				errorKey:    keyEmail,
				wantMessage: MsgEmailRequired,
			},
			{
				name:        "wrong email format",
				rawPayload:  `{"email":"test@com" }`,
				errorKey:    keyEmail,
				wantMessage: MsgWrongEmailFormat,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := newPOSTUserHTTPRequest(tc.rawPayload)
				res := httptest.NewRecorder()
				handler.ServeHTTP(res, req)

				var got api.ErrorResponse
				assert.Equal(t, res.Code, http.StatusUnprocessableEntity, "did not get correct response status code")
				assert.Equal(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
				assert.JSONDecode(t, res.Body, &got)
				assert.Equal(t, got.Message, api.MsgInvalidReqPayload, "wrong error Message in api.ErrorResponse")

				msg, exists := got.Errors[tc.errorKey]
				assert.Equal(t, exists, true, "expected key in api.ErrorResponse.Errors map")
				assert.Equal(t, msg, tc.wantMessage, "wrong error message in api.ErrorResponse.Errors")
			})
		}

		t.Run("missing name AND email", func(t *testing.T) {
			rawPayload := `{}`
			req := newPOSTUserHTTPRequest(rawPayload)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var got api.ErrorResponse

			assert.Equal(t, res.Code, http.StatusUnprocessableEntity, "did not get correct response status code")

			assert.Equal(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
			assert.JSONDecode(t, res.Body, &got)
			assert.Equal(t, got.Message, api.MsgInvalidReqPayload, "wrong 'Message' field in api.ErrorResponse")

			nameMsg, hasName := got.Errors[keyName]
			assert.Equal(t, hasName, true, "expected 'name' key in api.ErrorResponse.Errors map")
			assert.Equal(t, nameMsg, MsgNameRequired, "wrong error message for 'name' field in api.ErrorResponse.Errors")

			emailMsg, hasEmail := got.Errors[keyEmail]
			assert.Equal(t, hasEmail, true, "expected 'email' key in api.ErrorResponse.Errors map")
			assert.Equal(t, emailMsg, MsgEmailRequired, "wrong error message in api.ErrorResponse.Error")
		})
	})
}

func TestPUT(t *testing.T) {
	testUserID := "cb392-zbfdau-adufba-audbfuda"

	t.Run("succesfully PUT user Name", func(t *testing.T) {
		repo := newStubRepository(testUserID)
		service := NewService(repo)
		handler := NewHandler(service)

		rawPayload := `{"name":"joseph"}`
		req := newPUTUserHTTPRequest(testUserID, rawPayload)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got User
		var expected UpdateUserRequest

		assert.Equal(t, res.Code, http.StatusOK, "did not get correct response status code")
		assert.Equal(t, res.Header().Get("content-type"), "application/json", "wrong content type format")

		t.Run("response body is correct", func(t *testing.T) {
			assert.JSONDecode(t, res.Body, &got)
			assert.Unmarshal(t, rawPayload, &expected)
			assert.Equal(t, got.Name, expected.Name, "user Name was not correctly updated in response body")
		})
		t.Run("updated user is correctly persisted in repository", func(t *testing.T) {
			updatedUser := repo.users[testUserID]
			assert.Equal(t, updatedUser.Name, expected.Name, "user Name was not correctly persisted in repository")
		})
	})

	t.Run("covers error scenarios", func(t *testing.T) {
		t.Run("wrong ID", func(t *testing.T) {
			fakeID := "aeh23-sf456-ae244-sf4"
			repo := newStubRepository(testUserID)
			service := NewService(repo)
			handler := NewHandler(service)

			rawPayload := `{"name":"joseph"}`
			req := newPUTUserHTTPRequest(fakeID, rawPayload)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var got api.ErrorResponse

			assert.Equal(t, res.Code, http.StatusNotFound, "did not get correct response status code")
			assert.Equal(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
			assert.JSONDecode(t, res.Body, &got)
			assert.Equal(t, got.Message, MsgNotFound, "wrong Message in ErrorResponse")
		})

		t.Run("missing Name", func(t *testing.T) {
			repo := newStubRepository(testUserID)
			service := NewService(repo)
			handler := NewHandler(service)

			rawPayload := `{}`
			req := newPUTUserHTTPRequest(testUserID, rawPayload)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var got api.ErrorResponse

			assert.Equal(t, res.Code, http.StatusUnprocessableEntity, "did not get correct response status code")
			assert.Equal(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
			assert.JSONDecode(t, res.Body, &got)
			assert.Equal(t, got.Message, api.MsgInvalidReqPayload, "wrong 'Message' field in api.ErrorResponse")
			nameMsg, hasName := got.Errors[keyName]
			assert.Equal(t, hasName, true, "expected 'name' key in api.ErrorResponse.Errors map")
			assert.Equal(t, nameMsg, MsgNameRequired, "wrong error message for 'name' field in api.ErrorResponse.Errors")
		})
	})

}

func newGETUserHTTPRequest(userID string) *http.Request {
	path := fmt.Sprintf("/users/%s", userID)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue("id", userID)
	return req
}

func newPOSTUserHTTPRequest(rawPayload string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(rawPayload))
	return req
}

func newPUTUserHTTPRequest(userID, rawPayload string) *http.Request {
	path := fmt.Sprintf("/users/%s", userID)
	req := httptest.NewRequest(http.MethodPut, path, strings.NewReader(rawPayload))
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

func assertUUIDFormat(t *testing.T, id string) {
	t.Helper()

	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("user ID is not in correct UUID format: %v", id)
	}
}

func assertTimeFormat(t *testing.T, time time.Time) {
	if time.IsZero() {
		t.Errorf("did not get correct CreatedAt in response body. It should not be zero")
	}
}
