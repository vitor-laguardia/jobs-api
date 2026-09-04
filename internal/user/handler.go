package user

import (
	"encoding/json"
	"net/http"

	"github.com/vitor-laguardia/jobs-api/internal/shared/httputil"
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
	user, err := h.service.GetByID(userID)

	if err != nil {
		errRes := httputil.NewErrorResponse(err.Error(), http.StatusNotFound, nil)
		w.WriteHeader(errRes.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
