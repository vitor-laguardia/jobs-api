package main

import (
	"log"
	"net/http"

	"github.com/vitor-laguardia/jobs-api/internal/job"
	"github.com/vitor-laguardia/jobs-api/internal/user"
)

func NewServer(jobHandler *job.Handler, userHandler *user.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/jobs/", jobHandler)
	mux.Handle("/jobs", jobHandler)
	mux.Handle("/users/", userHandler)
	mux.Handle("/users", userHandler)
	return mux
}

func main() {
	jobRepo := job.NewInMemoryRepository()
	jobService := job.NewService(jobRepo)
	jobHandler := job.NewHandler(jobService)

	userRepo := user.NewInMemoryRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	server := NewServer(jobHandler, userHandler)

	log.Fatal(http.ListenAndServe(":8080", server))
}
