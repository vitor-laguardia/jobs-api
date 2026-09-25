package job

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/shared/assert"
	"github.com/vitor-laguardia/jobs-api/internal/user"
)

func TestPOSTJobAndGETIt(t *testing.T) {
	userRepo := user.NewInMemoryRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	jobRepo := NewInMemoryRepository(userRepo)
	jobService := NewService(jobRepo)
	jobHandler := NewHandler(jobService)

	userPayload := strings.NewReader(`{"name":"john", "email":"john@gmail.com"}`)
	userReq := httptest.NewRequest(http.MethodPost, "/users", userPayload)
	userRes := httptest.NewRecorder()
	userHandler.ServeHTTP(userRes, userReq)

	if userRes.Code != http.StatusCreated {
		t.Fatalf("setup failed: could not create user, status: %d body: %s", userRes.Code, userRes.Body.String())
	}

	var createdUser user.User
	if err := json.NewDecoder(userRes.Body).Decode(&createdUser); err != nil {
		t.Fatalf("setup failed: could not decode user response: %v", err)
	}

	jobPayload := fmt.Sprintf(
		`{"title":"my first job","description":"getting used to it","priority":1,"userId":%q}`,
		createdUser.ID,
	)

	POSTReq := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(jobPayload))
	POSTRes := httptest.NewRecorder()
	jobHandler.ServeHTTP(POSTRes, POSTReq)

	if POSTRes.Code != http.StatusCreated {
		t.Fatalf("POST /jobs failed, status: %d body: %s", POSTRes.Code, POSTRes.Body.String())
	}

	var createdJob Job
	if err := json.NewDecoder(POSTRes.Body).Decode(&createdJob); err != nil {
		t.Fatalf("could not decode job response: %v", err)
	}

	GETReq := httptest.NewRequest(http.MethodGet, "/jobs/"+createdJob.ID, nil)
	GETReq.SetPathValue("id", createdJob.ID)
	GETRes := httptest.NewRecorder()
	jobHandler.ServeHTTP(GETRes, GETReq)

	assert.Equal(t, GETRes.Code, http.StatusOK, "wrong response status (GET job)")

	var got Job
	assert.JSONDecode(t, GETRes.Body, &got)
	assertJob(t, got, createdJob)
}
