package main

import (
	"log"
	"net/http"

	"github.com/vitor-laguardia/jobs-api/internal/job"
)

func main() {
	jobs := job.NewInMemoryJobs()
	service := job.NewJobService(jobs)
	jh := job.NewJobHandler(service)

	log.Fatal(http.ListenAndServe(":8080", jh))
}
