package user

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
	router  *http.ServeMux
}

func NewHandler(service *Service) *Handler {
	h := &Handler{service: service}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", h.getUser)
	h.router = mux
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	user := h.service.GetByID(userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
