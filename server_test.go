package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubJobService struct {
	jobs map[string]string
}

func (s *StubJobService) GetById(jobId string) string {
	job := s.jobs[jobId]
	return job
}

func TestGETJob(t *testing.T) {
	jobs := &StubJobService{
		map[string]string{
			"1": "job1", "2": "job2"},
	}
	server := &JobServer{job: jobs}

	t.Run("return job 1 information", func(t *testing.T) {
		req := newGetJobsRequest("1")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusOK)
		assertResponseBody(t, res.Body.String(), "job1")
	})

	t.Run("return job 2 information", func(t *testing.T) {
		req := newGetJobsRequest("2")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusOK)
		assertResponseBody(t, res.Body.String(), "job2")
	})

	t.Run("returns 404 on missing jobs", func(t *testing.T) {
		req := newGetJobsRequest("3")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		assertStatus(t, res.Code, http.StatusNotFound)
	})

}

func newGetJobsRequest(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/jobs/%s", id), nil)
	return req
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q, want %q", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
