package job

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vitor-laguardia/jobs-api/internal/shared/assert"
)

func TestPOSTJobAndGETIt(t *testing.T) {
	jobDb := NewInMemoryJobs()
	service := NewJobService(jobDb)
	server := NewJobHandler(service)

	rawPayload := `{"title":"my first job","description":"getting use with it","priority":1,"userId":"1"}`
	POSTReq := newPOSTJobHTTPRequest(rawPayload)
	POSTRes := httptest.NewRecorder()

	server.ServeHTTP(POSTRes, POSTReq)

	var newJob Job
	json.NewDecoder(POSTRes.Body).Decode(&newJob)

	path := fmt.Sprintf("/jobs/%s", newJob.ID)
	GETReq := httptest.NewRequest(http.MethodGet, path, nil)
	GETRes := httptest.NewRecorder()

	GETReq.SetPathValue("id", newJob.ID)
	server.ServeHTTP(GETRes, GETReq)

	var got Job

	assert.Equal(t, GETRes.Code, http.StatusOK, "wrong response status (GET Job): ")
	assert.JSONDecode(t, GETRes.Body, &got)
	assert.Equal(t, GETRes.Code, http.StatusOK, "did not get correct status")
	assertJob(t, got, newJob)
}
