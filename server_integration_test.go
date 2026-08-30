package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPOSTJobAndGETIt(t *testing.T) {
	jobDb := NewInMemoryJobs()
	service := newJobService(jobDb)
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

	assertEqual(t, GETRes.Code, http.StatusOK, "wrong response status (GET Job): ")
	assertJSONDecode(t, GETRes.Body, &got)
	assertEqual(t, GETRes.Code, http.StatusOK, "did not get correct status")
	assertJob(t, got, newJob)
}
