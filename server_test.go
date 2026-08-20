package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type StubJobRepository struct {
	jobs map[string]*Job
}

func (s *StubJobRepository) GetByID(jobID string) *Job {
	job, ok := s.jobs[jobID]

	if !ok {
		return nil
	}

	jobCopy := *job
	return &jobCopy
}

func (s *StubJobRepository) Create(job *Job) {
	s.jobs[job.ID] = job
}

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated network read failure")
}

func TestGETJob(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	jobStub := &StubJobRepository{
		map[string]*Job{
			"1": &Job{
				ID:          uuid.New().String(),
				Title:       "GET job title",
				Description: "lorem ipsum",
				Status:      JobStatusPending,
				Priority:    JobPriority(2),
				UserID:      "1",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}
	server := &JobServer{repo: jobStub}

	t.Run("get job stub", func(t *testing.T) {
		req := newGETJobHTTPRequest("1")
		res := httptest.NewRecorder()

		req.SetPathValue("id", "1")
		server.ServeHTTP(res, req)

		var got Job
		want := *jobStub.jobs["1"]

		assertEqual(t, res.Code, http.StatusOK, "did not get correct status")
		assertEqual(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
		assertJSONDecode(t, res.Body, &got)
		assertJob(t, got, want)
	})

	t.Run("returns 404 on missing jobs", func(t *testing.T) {
		req := newGETJobHTTPRequest("3")
		res := httptest.NewRecorder()

		req.SetPathValue("id", "3")
		server.ServeHTTP(res, req)

		var got ErrorResponse

		assertEqual(t, res.Code, http.StatusNotFound, "did not get correct status")
		assertEqual(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
		assertJSONDecode(t, res.Body, &got)
		assertEqual(t, got.Message, MsgJobNotFound, "wrong Message in error response")
	})

}

func TestPostJob(t *testing.T) {
	stubRepo := &StubJobRepository{make(map[string]*Job)}
	server := &JobServer{repo: stubRepo}

	t.Run("it returns the new job as JSON", func(t *testing.T) {
		rawPayload := `{"title": "POST job", "description": "Job for testing purpose", "priority": 2, "userId": "1"}`

		req := newPOSTJobHTTPRequest(rawPayload)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got Job
		var want CreateJobRequest

		assertEqual(t, res.Code, http.StatusCreated, "did not get correct status")
		assertEqual(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
		assertJSONDecode(t, res.Body, &got)
		assertUnmarshal(t, rawPayload, &want)
		assertPostJobResponseBody(t, got, want)
	})

	validationTests := []struct {
		name        string
		payload     string
		wantStatus  int
		wantMessage string
		wantKey     string
	}{
		{
			name:        "missing title",
			payload:     `{"title": "", "description": "Job for testing purpose", "priority": 2, "userId": "1"}`,
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: MsgTitleRequired,
			wantKey:     KeyTitle,
		},
		{
			name:        "missing userId",
			payload:     `{"title": "POST job", "description": "Job for testing purpose", "priority": 2, "userId": ""}`,
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: MsgUserIDRequired,
			wantKey:     KeyUserID,
		},
		{
			name:        "priority in negative range",
			payload:     `{"title": "POST job", "description": "Job for testing purpose", "priority": -1, "userId": "1"}`,
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: MsgPriorityOutOfRange,
			wantKey:     KeyPriority,
		},
		{
			name:        "priority out of range",
			payload:     `{"title": "POST job", "description": "Job for testing purpose", "priority": 4, "userId": "1"}`,
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: MsgPriorityOutOfRange,
			wantKey:     KeyPriority,
		},
	}

	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			req := newPOSTJobHTTPRequest(tt.payload)
			res := httptest.NewRecorder()

			server.ServeHTTP(res, req)

			var got ErrorResponse
			assertEqual(t, res.Code, tt.wantStatus, "did not get correct status")
			assertEqual(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
			assertJSONDecode(t, res.Body, &got)
			assertEqual(t, got.Message, MsgInvalidReqPayload, "wrong message in error response")
			for key, errMsg := range got.Errors {
				assertEqual(t, key, tt.wantKey, "wrong Key in error response")
				assertEqual(t, errMsg, tt.wantMessage, "wrong error message in response body")
			}
		})
	}

}

func TestPayload(t *testing.T) {
	stubRepo := &StubJobRepository{make(map[string]*Job)}
	server := &JobServer{repo: stubRepo}

	tableTests := []struct {
		name        string
		body        string
		needSprintf bool
		formatValue any
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "request with empty body",
			body:        ``,
			needSprintf: false,
			formatValue: 0,
			wantStatus:  http.StatusBadRequest,
			wantMessage: MsgEmptyBody,
		},
		{
			name:        "request with unclosed json",
			body:        `{"title": "testing"`,
			needSprintf: false,
			formatValue: 0,
			wantStatus:  http.StatusBadRequest,
			wantMessage: MsgUnexpectedEOF,
		},
		{
			name:        "request with syntax error",
			body:        `{"title": testing`,
			needSprintf: true,
			formatValue: 12,
			wantStatus:  http.StatusBadRequest,
			wantMessage: MsgSyntaxErr,
		},
		{
			name:        "request with unknown json field",
			body:        `{"title": "testing", "badField": 156}`,
			needSprintf: true,
			formatValue: `"badField"`,
			wantStatus:  http.StatusBadRequest,
			wantMessage: MsgUnknownField,
		},
		{
			name:        "request with more then one body",
			body:        `{"title": "testing", "userId": "2"}{"otherBody": "hack"}`,
			needSprintf: false,
			formatValue: 0,
			wantStatus:  http.StatusBadRequest,
			wantMessage: MsgMultipleReqBody,
		},
	}

	for _, tt := range tableTests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(tt.body))
			res := httptest.NewRecorder()

			server.ServeHTTP(res, req)

			var got ErrorResponse
			assertJSONDecode(t, res.Body, &got)

			wantMessage := func() string {
				if tt.needSprintf {
					return fmt.Sprintf(tt.wantMessage, tt.formatValue)
				}
				return tt.wantMessage
			}

			assertEqual(t, res.Code, tt.wantStatus, "did not get correct status")
			assertEqual(t, got.Message, wantMessage(), "wrong error message in response body")
		})
	}

	t.Run("request with invalid field type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(`{"title": 123}`))
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got ErrorResponse
		want := fmt.Sprintf(MsgWrongFieldType, "title", 13)

		assertEqual(t, res.Code, http.StatusBadRequest, "did not get correct status")
		assertJSONDecode(t, res.Body, &got)
		assertEqual(t, got.Message, want, "wrong error message in response body")
	})

	t.Run("request body with 1MB", func(t *testing.T) {
		largePadding := strings.Repeat("a", MaxBodyBytes+1)
		largeJSON := fmt.Sprintf(`{"title": "%s"}`, largePadding)
		limitExceededBody := strings.NewReader(largeJSON)
		req := httptest.NewRequest(http.MethodPost, "/jobs", limitExceededBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got ErrorResponse
		want := fmt.Sprintf(MsgMaxBodyBytes, MaxBodyBytes)

		assertJSONDecode(t, res.Body, &got)

		assertEqual(t, res.Code, http.StatusRequestEntityTooLarge, "did not get correct status")
		assertEqual(t, got.Message, want, "wrong error message in response body")
	})

	t.Run("request header without application/json", func(t *testing.T) {
		rawTextBody := "id=123&status=active"
		req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(rawTextBody))
		req.Header.Set("content-type", " text/plain")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got ErrorResponse
		assertJSONDecode(t, res.Body, &got)

		assertEqual(t, res.Code, http.StatusUnsupportedMediaType, "did not get correct status")
		assertEqual(t, got.Message, MsgUnexpectedContentType, "wrong error message in response body")

	})

	t.Run("unmaped JSON error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/jobs", errReader{})
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got ErrorResponse

		assertJSONDecode(t, res.Body, &got)
		assertEqual(t, res.Code, http.StatusInternalServerError, "did not get correct status")
		assertEqual(t, got.Message, MsgJSONError, "wrong error message in response body")
	})
}

func newGETJobHTTPRequest(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/jobs/%s", id), nil)
	return req
}

func newPOSTJobHTTPRequest(rawPayload string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(rawPayload))
	return req
}

func assertPostJobResponseBody(t *testing.T, got Job, want CreateJobRequest) {
	t.Helper()

	assertEqual(t, got.Title, want.Title, "response body is wrong")
	assertEqual(t, got.Description, want.Description, "response body is wrong")
	assertEqual(t, got.Priority, JobPriority(want.Priority), "response body is wrong")
	assertEqual(t, got.UserID, want.UserID, "response body is wrong")
	if _, err := uuid.Parse(got.ID); err != nil {
		t.Errorf("response body is wrong,  not UUID format: %v", got.ID)
	}
	if got.CreatedAt.IsZero() {
		t.Errorf("response body is wrong, CreatedAt should not be zero")
	}
}
