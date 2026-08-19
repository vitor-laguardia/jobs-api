package main

import (
	"log"
	"net/http"
)

func main() {
	jobs := NewInMemoryJobs()
	server := &JobServer{jobs}
	http.Handle("/jobs/{id}", server)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
