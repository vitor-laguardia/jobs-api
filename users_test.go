package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETUser(t *testing.T) {
	testUserID := "ab342-sbfdau-adufba-audbfuda"
	t.Run("returns 200 to get user", func(t *testing.T) {
		repo := setupRepo(testUserID)
		service := newJobService(repo)
		jobHandler := NewJobHandler(service)

		req := httptest.NewRequest(http.MethodGet, "/user", nil)
		res := httptest.NewRecorder()

		jobHandler.ServeHTTP(res, req)

		assertEqual(t, res.Code, http.StatusOK, "did not get correct response status code")

	})
}
