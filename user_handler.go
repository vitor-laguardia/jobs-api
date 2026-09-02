package main

import (
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	service *UserService
	router  *http.ServeMux
}

func NewUserHandler() *UserHandler {
	uh := &UserHandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", uh.getUser)
	uh.router = mux
	return uh
}

func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uh.router.ServeHTTP(w, r)
}

func (uh *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	user := uh.service.GetByID(userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
