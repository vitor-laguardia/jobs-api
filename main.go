package main

import (
	"log"
	"net/http"
)

type InMemoryJobs struct{}

func (i *InMemoryJobs) GetByID(jobId string) string {
	return "123"
}

func main() {
	jobs := &InMemoryJobs{}
	server := &JobServer{jobs}
	http.Handle("/jobs/{id}", server)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
