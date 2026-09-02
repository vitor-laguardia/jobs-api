package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETUser(t *testing.T) {
	testUserID := "ab342-sbfdau-adufba-audbfuda"
	t.Run("returns 200 to get user", func(t *testing.T) {
		userHandler := NewUserHandler()

		req := newGETUserHTTPRequest(testUserID)
		res := httptest.NewRecorder()

		userHandler.ServeHTTP(res, req)

		var got User

		assertEqual(t, res.Code, http.StatusOK, "did not get correct response status code")
		assertJSONDecode(t, res.Body, &got)
		assertEqual(t, got.Name, "alfred", "did not get correct user Name")
		assertEqual(t, got.ID, testUserID, "did not get correct user ID")
		assertEqual(t, got.Email, "alfred@gmail.com", "did not get correct user email")
	})
}

func newGETUserHTTPRequest(userID string) *http.Request {
	path := fmt.Sprintf("/users/%s", userID)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue("id", userID)
	return req
}
