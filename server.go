package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	MsgInvalidReqPayload = "invalid request payload"
)

type JobService interface {
	GetByID(id string) string
}

type JobServer struct {
	job JobService
}

type Validator interface {
	Valid() (problems map[string]string)
}

type ErrorResponse struct {
	Message string            `json:"message"`          // Resumo legível do erro
	Errors  map[string]string `json:"errors,omitempty"` // Detalhes por campo (se houver)
}

func (j *JobServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		j.getJob(w, r)

	case http.MethodPost:
		j.postJob(w, r)
	}
}

func (j *JobServer) postJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jobReq, problems, err := decodeValid[CreateJobRequest](r)

	//TODO: include marshal error validation in tests
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errRes := ErrorResponse{
			Message: MsgInvalidReqPayload,
			Errors:  problems,
		}
		json.NewEncoder(w).Encode(errRes)
		return
	}

	newJob := NewJob(jobReq.Title, jobReq.Description, JobPriority(jobReq.Priority), jobReq.UserID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newJob)
}

func (j *JobServer) getJob(w http.ResponseWriter, r *http.Request) {
	jobId := strings.TrimPrefix(r.URL.Path, "/jobs/")
	job := j.job.GetByID(jobId)

	if job == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, job)
}

func decodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}
	if problems := v.Valid(); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}
