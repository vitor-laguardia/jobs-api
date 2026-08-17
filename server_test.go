package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type StubJobService struct {
	jobs map[string]*Job
}

func (s *StubJobService) GetByID(jobID string) *Job {
	job, ok := s.jobs[jobID]

	if !ok {
		return nil
	}

	jobCopy := *job
	return &jobCopy
}

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated network read failure")
}
func TestGETJob(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	jobStub := &StubJobService{
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
	server := &JobServer{job: jobStub}

	t.Run("get job stub", func(t *testing.T) {
		req := newGETJobHTTPRequest("1")
		res := httptest.NewRecorder()

		req.SetPathValue("id", "1")
		server.ServeHTTP(res, req)

		var got Job
		want := *jobStub.jobs["1"]

		assertStatus(t, res.Code, http.StatusOK)
		assertHeader(t, res.Header().Get("content-type"), "application/json")
		assertJSONDecode(t, res.Body, &got)
		assertJob(t, got, want)
	})

	t.Run("returns 404 on missing jobs", func(t *testing.T) {
		req := newGETJobHTTPRequest("3")
		res := httptest.NewRecorder()

		req.SetPathValue("id", "3")
		server.ServeHTTP(res, req)

		var got ErrorResponse

		assertStatus(t, res.Code, http.StatusNotFound)
		assertHeader(t, res.Header().Get("content-type"), "application/json")
		assertJSONDecode(t, res.Body, &got)
		assertString(t, got.Message, MsgJobNotFound, "wrong Message in error response")
	})

}

func TestPostJob(t *testing.T) {
	jobs := &StubJobService{nil}
	server := &JobServer{job: jobs}

	t.Run("it returns the new job as JSON", func(t *testing.T) {
		payload := CreateJobRequest{
			Title:       "POST job",
			Description: "Job for testing purpose",
			Priority:    2,
			UserID:      "1"}

		req := newPOSTJobHTTPRequest(t, payload)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got Job

		assertStatus(t, res.Code, http.StatusCreated)
		assertHeader(t, res.Header().Get("content-type"), "application/json")
		assertJSONDecode(t, res.Body, &got)
		assertPostJobResponseBody(t, got, payload)
	})

	validationTests := []struct {
		name        string
		payload     CreateJobRequest
		wantStatus  int
		wantMessage string
		wantKey     string
	}{
		{
			name: "missing title",
			payload: CreateJobRequest{
				Title:       "",
				Description: "Job for testing",
				Priority:    2,
				UserID:      "user-1",
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: MsgTitleRequired,
			wantKey:     KeyTitle,
		},
		{
			name: "missing userId",
			payload: CreateJobRequest{
				Title:       "POST Job",
				Description: "Job for testing",
				Priority:    2,
				UserID:      "",
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: MsgUserIDRequired,
			wantKey:     KeyUserID,
		},
		{
			name: "priority in negative range",
			payload: CreateJobRequest{
				Title:       "POST Job",
				Description: "Job for testing",
				Priority:    -1,
				UserID:      "user-1",
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: MsgPriorityOutOfRange,
			wantKey:     KeyPriority,
		},
		{
			name: "priority out of range",
			payload: CreateJobRequest{
				Title:       "POST Job",
				Description: "Job for testing",
				Priority:    4,
				UserID:      "user-1",
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: MsgPriorityOutOfRange,
			wantKey:     KeyPriority,
		},
	}

	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			req := newPOSTJobHTTPRequest(t, tt.payload)
			res := httptest.NewRecorder()

			server.ServeHTTP(res, req)

			var got ErrorResponse
			assertStatus(t, res.Code, tt.wantStatus)
			assertHeader(t, res.Header().Get("content-type"), "application/json")
			assertJSONDecode(t, res.Body, &got)
			assertString(t, got.Message, MsgInvalidReqPayload, "wrong message in error response")
			for key, errMsg := range got.Errors {
				assertString(t, key, tt.wantKey, "wrong Key in error response")
				assertString(t, errMsg, tt.wantMessage, "wrong error message in response body")
			}
		})
	}

}

func TestPayload(t *testing.T) {
	jobs := &StubJobService{nil}
	server := &JobServer{job: jobs}

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

			assertStatus(t, res.Code, tt.wantStatus)
			assertString(t, got.Message, wantMessage(), "wrong error message in response body")
		})
	}

	t.Run("request with invalid field type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(`{"title": 123}`))
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got ErrorResponse
		want := fmt.Sprintf(MsgWrongFieldType, "title", 13)

		assertJSONDecode(t, res.Body, &got)
		assertStatus(t, res.Code, http.StatusBadRequest)
		assertString(t, got.Message, want, "wrong error message in response body")
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
		assertStatus(t, res.Code, http.StatusRequestEntityTooLarge)
		assertString(t, got.Message, want, "wrong error message in response body")
	})

	t.Run("request header without application/json", func(t *testing.T) {
		rawTextBody := "id=123&status=active"
		req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(rawTextBody))
		req.Header.Set("content-type", " text/plain")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got ErrorResponse
		assertJSONDecode(t, res.Body, &got)
		assertStatus(t, res.Code, http.StatusUnsupportedMediaType)
		assertString(t, got.Message, MsgUnexpectedContentType, "wrong error message in response body")

	})

	t.Run("default case - generic read error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/jobs", errReader{})
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got ErrorResponse

		assertJSONDecode(t, res.Body, &got)
		assertStatus(t, res.Code, http.StatusInternalServerError)
		assertString(t, got.Message, MsgJSONError, "wrong error message in response body")
	})
}

func newGETJobHTTPRequest(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/jobs/%s", id), nil)
	return req
}

func newPOSTJobHTTPRequest(t *testing.T, payload CreateJobRequest) *http.Request {
	t.Helper()

	bodyBytes, err := json.Marshal(payload)

	if err != nil {
		t.Fatalf("Fail in serialize POST job payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(bodyBytes))
	return req
}

func assertHeader(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Wrong content type format. Want: %q, got: %q", want, got)
	}
}

func assertJSONDecode(t *testing.T, r io.Reader, v interface{}) {
	t.Helper()
	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		t.Fatalf("failed to decode JSON: %v in Response: %v", err, r)
	}
}

func assertString(t *testing.T, got, want, message string) {
	t.Helper()

	if got != want {
		if message != "" {
			t.Errorf("%s: got %q, want %q", message, got, want)
			return
		}
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertPostJobResponseBody(t *testing.T, got Job, want CreateJobRequest) {
	t.Helper()
	if got.Title != want.Title {
		t.Errorf("response body is wrong, got Title %q, want %q", got.Title, want.Title)
	}
	if got.Description != want.Description {
		t.Errorf("response body is wrong, got Description  %q, want %q", got.Description, want.Description)
	}
	if int(got.Priority) != want.Priority {
		t.Errorf("response body is wrong, got Priority  %d, want %d", got.Priority, want.Priority)
	}
	if got.UserID != want.UserID {
		t.Errorf("response body is wrong, got UserID  %q, want %q", got.UserID, want.UserID)
	}
	if _, err := uuid.Parse(got.ID); err != nil {
		t.Errorf("response body is wrong,  not UUID format: %v", got.ID)
	}
	if got.CreatedAt.IsZero() {
		t.Errorf("response body is wrong, CreatedAt should not be zero")
	}
}

func assertJob(t *testing.T, got Job, want Job) {
	t.Helper()
	assertString(t, got.ID, want.ID, "wrong id in response body")
	assertString(t, got.Title, want.Title, "wrong title in response body")
	assertString(t, got.Description, want.Description, "wrong description in response body")
	assertString(t, got.UserID, want.UserID, "wrong userID in response body")
	if got.Priority != want.Priority {
		t.Errorf("wrong Priority in response body: got %d, want %d", got.Priority, want.Priority)
	}
	if got.Status != want.Status {
		t.Errorf("wrong Status in response body: got %q, want %q", got.Status, want.Status)
	}
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
