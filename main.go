package main

import (
	"log"
	"net/http"
)

func main() {
	jobs := NewInMemoryJobs()
	service := newJobService(jobs)
	jh := NewJobHandler(service)

	log.Fatal(http.ListenAndServe(":8080", jh))
}
