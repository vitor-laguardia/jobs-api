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

func (s *StubJobRepository) Update(job *Job) *Job {
	if _, exists := s.jobs[job.ID]; !exists {
		return nil
	}
	s.jobs[job.ID] = job
	return job
}

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated network read failure")
}

func TestGETJob(t *testing.T) {
	const testJobID = "b2z4x3d4-e5f6-7890-1234-56789abcdef0"

	jobStub := setupRepo(testJobID)
	service := newJobService(jobStub)
	jobHandler := NewJobHandler(service)

	t.Run("get job stub", func(t *testing.T) {
		req := newGETJobHTTPRequest(testJobID)
		res := httptest.NewRecorder()

		req.SetPathValue("id", testJobID)
		jobHandler.ServeHTTP(res, req)

		var got Job
		want := *jobStub.jobs[testJobID]

		assertEqual(t, res.Code, http.StatusOK, "did not get correct response status code")
		assertEqual(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
		assertJSONDecode(t, res.Body, &got)
		assertJob(t, got, want)
	})

	t.Run("returns 404 on missing jobs", func(t *testing.T) {
		fakeJobID := "w4z4x3r4-e5f6-7890-1234-56789abcdef0"
		req := newGETJobHTTPRequest(fakeJobID)
		res := httptest.NewRecorder()

		req.SetPathValue("id", fakeJobID)
		jobHandler.ServeHTTP(res, req)

		var got ErrorResponse

		assertEqual(t, res.Code, http.StatusNotFound, "did not get correct response status code")
		assertEqual(t, res.Header().Get("content-type"), "application/json", "wrong content type format")
		assertJSONDecode(t, res.Body, &got)
		assertEqual(t, got.Message, MsgJobNotFound, "wrong Message in error response")
	})

}

func TestPostJob(t *testing.T) {
	stubRepo := &StubJobRepository{make(map[string]*Job)}
	s := newJobService(stubRepo)
	jh := NewJobHandler(s)

	t.Run("it returns the new job as JSON", func(t *testing.T) {
		rawPayload := `{"title": "POST job", "description": "Job for testing purpose", "priority": 2, "userId": "1"}`

		req := newPOSTJobHTTPRequest(rawPayload)
		res := httptest.NewRecorder()

		jh.ServeHTTP(res, req)

		var got Job
		var want CreateJobRequest

		assertEqual(t, res.Code, http.StatusCreated, "did not get correct response status code")
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

			jh.ServeHTTP(res, req)

			var got ErrorResponse
			assertEqual(t, res.Code, tt.wantStatus, "did not get correct response status code")
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
	s := newJobService(stubRepo)
	jh := NewJobHandler(s)

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

			jh.ServeHTTP(res, req)

			var got ErrorResponse
			assertJSONDecode(t, res.Body, &got)

			wantMessage := func() string {
				if tt.needSprintf {
					return fmt.Sprintf(tt.wantMessage, tt.formatValue)
				}
				return tt.wantMessage
			}

			assertEqual(t, res.Code, tt.wantStatus, "did not get correct response status code")
			assertEqual(t, got.Message, wantMessage(), "wrong error message in response body")
		})
	}

	t.Run("request with invalid field type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(`{"title": 123}`))
		res := httptest.NewRecorder()

		jh.ServeHTTP(res, req)

		var got ErrorResponse
		want := fmt.Sprintf(MsgWrongFieldType, "title", 13)

		assertEqual(t, res.Code, http.StatusBadRequest, "did not get correct response status code")
		assertJSONDecode(t, res.Body, &got)
		assertEqual(t, got.Message, want, "wrong error message in response body")
	})

	t.Run("request body with 1MB", func(t *testing.T) {
		largePadding := strings.Repeat("a", MaxBodyBytes+1)
		largeJSON := fmt.Sprintf(`{"title": "%s"}`, largePadding)
		limitExceededBody := strings.NewReader(largeJSON)
		req := httptest.NewRequest(http.MethodPost, "/jobs", limitExceededBody)
		res := httptest.NewRecorder()

		jh.ServeHTTP(res, req)

		var got ErrorResponse
		want := fmt.Sprintf(MsgMaxBodyBytes, MaxBodyBytes)

		assertJSONDecode(t, res.Body, &got)

		assertEqual(t, res.Code, http.StatusRequestEntityTooLarge, "did not get correct response status code")
		assertEqual(t, got.Message, want, "wrong error message in response body")
	})

	t.Run("request header without application/json", func(t *testing.T) {
		rawTextBody := "id=123&status=active"
		req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(rawTextBody))
		req.Header.Set("content-type", " text/plain")
		res := httptest.NewRecorder()

		jh.ServeHTTP(res, req)

		var got ErrorResponse
		assertJSONDecode(t, res.Body, &got)

		assertEqual(t, res.Code, http.StatusUnsupportedMediaType, "did not get correct response status code")
		assertEqual(t, got.Message, MsgUnexpectedContentType, "wrong error message in response body")

	})

	t.Run("unmaped JSON error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/jobs", errReader{})
		res := httptest.NewRecorder()

		jh.ServeHTTP(res, req)

		var got ErrorResponse

		assertJSONDecode(t, res.Body, &got)
		assertEqual(t, res.Code, http.StatusInternalServerError, "did not get correct response status code")
		assertEqual(t, got.Message, MsgJSONError, "wrong error message in response body")
	})
}

func TestPUTJob(t *testing.T) {
	const testJobID = "a1b2c3d4-e5f6-7890-1234-56789abcdef0"

	testCases := []struct {
		name       string
		rawPayload string
		checkField func(t *testing.T, got, original Job)
	}{
		{
			name:       "updates only title",
			rawPayload: `{"title": "another title"}`,
			checkField: func(t *testing.T, got, original Job) {
				assertEqual(t, got.Title, "another title", "job title was not correctly updated")
				assertEqual(t, got.Description, original.Description, "job description was unexpectedly updated")
				assertEqual(t, got.Priority, original.Priority, "job priority was unexpectedly updated")
				assertEqual(t, got.Status, original.Status, "job status was unexpectedly updated")
			},
		},
		{
			name:       "updates only description",
			rawPayload: `{"description": "another description"}`,
			checkField: func(t *testing.T, got, original Job) {
				assertEqual(t, got.Description, "another description", "job description was not correctly updated")
				assertEqual(t, got.Title, original.Title, "job title was unexpectedly updated")
				assertEqual(t, got.Priority, original.Priority, "job priority was unexpectedly updated")
				assertEqual(t, got.Status, original.Status, "job status was unexpectedly updated")
			},
		},
		{
			name:       "updates only priority",
			rawPayload: `{"priority": 3}`,
			checkField: func(t *testing.T, got, original Job) {
				assertEqual(t, got.Priority, JobPriority(3), "job priority was not correctly updated")
				assertEqual(t, got.Title, original.Title, "job title was unexpectedly updated")
				assertEqual(t, got.Description, original.Description, "job description was unexpectedly updated")
				assertEqual(t, got.Status, original.Status, "job status was unexpectedly updated")
			},
		},
		{
			name:       "updates only status",
			rawPayload: `{"status": "running"}`,
			checkField: func(t *testing.T, got, original Job) {
				assertEqual(t, got.Status, JobStatusRunning, "job status was not correctly updated")
				assertEqual(t, got.Title, original.Title, "job title was unexpectedly updated")
				assertEqual(t, got.Description, original.Description, "job description was unexpectedly updated")
				assertEqual(t, got.Priority, original.Priority, "job priority was unexpectedly updated")
			},
		},
		{
			name:       "updates all fields together",
			rawPayload: `{"title": "another title", "description": "another description", "priority": 3, "status": "running"}`,
			checkField: func(t *testing.T, got, original Job) {
				assertEqual(t, got.Title, "another title", "job title was not correctly updated")
				assertEqual(t, got.Description, "another description", "job description was not correctly updated")
				assertEqual(t, got.Priority, JobPriority(3), "job priority was not correctly updated")
				assertEqual(t, got.Status, JobStatusRunning, "job status was not correctly updated")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stubRepo := setupRepo(testJobID)
			s := newJobService(stubRepo)
			jh := NewJobHandler(s)
			original := *stubRepo.jobs[testJobID]

			req := newPUTJobHTTPRequest(testJobID, tc.rawPayload)
			res := httptest.NewRecorder()

			jh.ServeHTTP(res, req)

			var got Job
			assertEqual(t, res.Code, http.StatusOK, "wrong response status")
			assertJSONDecode(t, res.Body, &got)
			assertTimeAfter(t, got.UpdatedAt, got.CreatedAt, "job UpdatedAt was not correctly updated")
			tc.checkField(t, got, original)
		})
	}

	t.Run("returns 422 on empty JSON body", func(t *testing.T) {
		stubRepo := &StubJobRepository{make(map[string]*Job)}
		s := newJobService(stubRepo)
		jh := NewJobHandler(s)

		rawPayload := `{}`
		req := newPUTJobHTTPRequest(testJobID, rawPayload)
		res := httptest.NewRecorder()

		jh.ServeHTTP(res, req)

		var got ErrorResponse
		assertJSONDecode(t, res.Body, &got)
		assertEqual(t, res.Code, http.StatusUnprocessableEntity, "did not get correct response status code")
		assertEqual(t, got.Message, MsgInvalidReqPayload, "wrong Message in ErrorResponse body")
		for key, errMsg := range got.Errors {
			assertEqual(t, key, KeyBody, "wrong Key in ErrorResponse.Errors")
			assertEqual(t, errMsg, MsgEmptyJSONBody, "wrong Message in ErrorResponse.Errors")
		}
	})

	t.Run("priority out of range", func(t *testing.T) {
		stubRepo := &StubJobRepository{make(map[string]*Job)}
		s := newJobService(stubRepo)
		jh := NewJobHandler(s)

		rawPayload := `{"priority": -1}`
		req := newPUTJobHTTPRequest(testJobID, rawPayload)
		res := httptest.NewRecorder()

		jh.ServeHTTP(res, req)

		var got ErrorResponse
		assertJSONDecode(t, res.Body, &got)

		assertEqual(t, res.Code, http.StatusUnprocessableEntity, "did not get correct response status code")
		assertEqual(t, got.Message, MsgInvalidReqPayload, "wrong Message in ErrorResponse body")
		for key, errMsg := range got.Errors {
			assertEqual(t, key, KeyPriority, "wrong Key in ErrorResponse.Errors")
			assertEqual(t, errMsg, MsgPriorityOutOfRange, "wrong message in ErrorResponse.Errors")
		}
	})

	t.Run("invalid status value", func(t *testing.T) {
		stubRepo := &StubJobRepository{make(map[string]*Job)}
		s := newJobService(stubRepo)
		jh := NewJobHandler(s)

		rawPayload := `{"status": "test"}`
		req := newPUTJobHTTPRequest(testJobID, rawPayload)
		res := httptest.NewRecorder()

		jh.ServeHTTP(res, req)

		var got ErrorResponse
		assertJSONDecode(t, res.Body, &got)

		assertEqual(t, res.Code, http.StatusUnprocessableEntity, "did not get correct response status code")
		assertEqual(t, got.Message, MsgInvalidReqPayload, "wrong Message in ErrorResponse body")
		for key, errMsg := range got.Errors {
			assertEqual(t, key, KeyStatus, "wrong Key in ErrorResponse.Errors")
			assertEqual(t, errMsg, MsgInvalidStatus, "wrong message in ErrorResponse.Errors")
		}
	})

	t.Run("invalid status transition", func(t *testing.T) {
		stubRepo := setupRepo(testJobID)
		s := newJobService(stubRepo)
		jh := NewJobHandler(s)

		rawPayload := `{"status": "done"}`
		req := newPUTJobHTTPRequest(testJobID, rawPayload)
		res := httptest.NewRecorder()

		jh.ServeHTTP(res, req)

		var got ErrorResponse
		assertJSONDecode(t, res.Body, &got)

		assertEqual(t, res.Code, http.StatusUnprocessableEntity, "did not get correct response status code")
		assertEqual(t, got.Message, MsgInvalidStatusTransition, "wrong Message in ErrorResponse body")
	})

	t.Run("update non-existent job", func(t *testing.T) {
		fakeJobID := "v2y4x5d4-e5f6-7890-1234-56789abcdef0"
		stubRepo := setupRepo(testJobID)
		service := newJobService(stubRepo)
		jobHandler := NewJobHandler(service)

		rawPayload := `{"title": "another title", "description": "another description", "priority": 3, "status": "running"}`
		req := newPUTJobHTTPRequest(fakeJobID, rawPayload)
		res := httptest.NewRecorder()

		jobHandler.ServeHTTP(res, req)

		var got ErrorResponse
		assertJSONDecode(t, res.Body, &got)
		assertEqual(t, res.Code, http.StatusBadRequest, "did not get correct response status code")
		assertEqual(t, got.Message, MsgJobNotFound, "wrong Message in ErrorResponse body")
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

func newPUTJobHTTPRequest(jobID, rawPayload string) *http.Request {
	path := fmt.Sprintf("/jobs/%s", jobID)
	req := httptest.NewRequest(http.MethodPut, path, strings.NewReader(rawPayload))
	req.SetPathValue("id", jobID)
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

func assertPutJobResponseBody(t *testing.T, got Job, want UpdateJobRequest) {
	t.Helper()
	fmt.Printf("\nINSIDE ASSERT job %#v \nwantDTO: %#v\n", got, want)

}

func setupRepo(testID string) *StubJobRepository {
	now := time.Now()
	jobStub := &StubJobRepository{
		map[string]*Job{
			testID: &Job{
				ID:          testID,
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
	return jobStub
}
