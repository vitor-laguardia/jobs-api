package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestPOSTJobAndGETIt(t *testing.T) {
	jobs := NewInMemoryJobs()
	server := JobServer{jobs}

	payloadString := `{"title":"my first job","description":"getting use with it","priority":1,"userId":"1"}`
	POSTReq := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(payloadString))
	POSTRes := httptest.NewRecorder()

	server.ServeHTTP(POSTRes, POSTReq)

	var job Job
	json.NewDecoder(POSTRes.Body).Decode(&job)
	fmt.Printf("\njob: %v", job)

	path := fmt.Sprintf("/jobs/%s", job.ID)
	GETReq := httptest.NewRequest(http.MethodGet, path, nil)
	GETRes := httptest.NewRecorder()

	GETReq.SetPathValue("id", job.ID)
	server.ServeHTTP(GETRes, GETReq)

	var got Job

	json.NewDecoder(GETRes.Body).Decode(&got)
	assertStatus(t, GETRes.Code, http.StatusOK)

	isEqual := reflect.DeepEqual(job, got)

	if !isEqual {
		t.Errorf("wrong integration between POST job and GET job: got %#v, want %#v", got, job)
	}
}
