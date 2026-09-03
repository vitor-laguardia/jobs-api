package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/shared/assert"
)

func TestGETUser(t *testing.T) {
	testUserID := "ab342-sbfdau-adufba-audbfuda"
	t.Run("returns 200 to get user", func(t *testing.T) {
		userHandler := NewUserHandler()

		req := newGETUserHTTPRequest(testUserID)
		res := httptest.NewRecorder()

		userHandler.ServeHTTP(res, req)

		var got User

		assert.Equal(t, res.Code, http.StatusOK, "did not get correct response status code")
		assert.JSONDecode(t, res.Body, &got)
		assert.Equal(t, got.Name, "alfred", "did not get correct user Name")
		assert.Equal(t, got.ID, testUserID, "did not get correct user ID")
		assert.Equal(t, got.Email, "alfred@gmail.com", "did not get correct user email")
	})
}

func newGETUserHTTPRequest(userID string) *http.Request {
	path := fmt.Sprintf("/users/%s", userID)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue("id", userID)
	return req
}
