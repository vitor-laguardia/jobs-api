package main

import (
	"fmt"
	"net/http"
	"strings"
)

type JobService interface {
	GetById(id string) string
}

type JobServer struct {
	job JobService
}

func (j *JobServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jobId := strings.TrimPrefix(r.URL.Path, "/jobs/")
	job := j.job.GetById(jobId)
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, job)
}
