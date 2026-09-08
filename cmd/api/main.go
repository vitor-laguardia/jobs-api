package main

import (
	"log"
	"net/http"

	"github.com/vitor-laguardia/jobs-api/internal/job"
	"github.com/vitor-laguardia/jobs-api/internal/user"
)

func NewServer(jobHandler *job.JobHandler, userHandler *user.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/jobs/", jobHandler)
	mux.Handle("/jobs", jobHandler)
	mux.Handle("/users/", userHandler)
	mux.Handle("/users", userHandler)
	return mux
}

func main() {
	jobRepo := job.NewInMemoryJobs()
	jobService := job.NewJobService(jobRepo)
	jobHandler := job.NewJobHandler(jobService)

	userRepo := user.NewInMemoryRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	server := NewServer(jobHandler, userHandler)

	log.Fatal(http.ListenAndServe(":8080", server))
}
