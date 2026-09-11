package user

import (
	"encoding/json"
	"net/http"

	"github.com/vitor-laguardia/jobs-api/internal/shared/api"
)

type Handler struct {
	service *Service
	router  *http.ServeMux
}

func NewHandler(service *Service) *Handler {
	h := &Handler{service: service}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", h.getUser)
	mux.HandleFunc("PUT /users/{id}", h.putUser)
	mux.HandleFunc("DELETE /users/{id}", h.deleteUser)
	mux.HandleFunc("POST /users", h.postUser)
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
		errRes := api.NewErrorResponse(err.Error(), http.StatusNotFound, nil)
		w.WriteHeader(errRes.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) postUser(w http.ResponseWriter, r *http.Request) {
	userInput, err := api.DecodeValid[CreateUserRequest](w, r)

	if err != nil {
		w.WriteHeader(err.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	newUser, serviceErr := h.service.Create(userInput)
	if serviceErr != nil {
		//TODO: tests dont cover yet because in memory repo dont have erros
		errRes := api.NewErrorResponse(serviceErr.Error(), http.StatusNotFound, nil)
		w.WriteHeader(errRes.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errRes)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func (h *Handler) putUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	userInput, decodeErr := api.DecodeValid[UpdateUserRequest](w, r)
	if decodeErr != nil {
		w.WriteHeader(decodeErr.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(decodeErr)
		return
	}

	newUser, updateErr := h.service.Update(userID, userInput)
	if updateErr != nil {
		errRes := api.NewErrorResponse(updateErr.Error(), http.StatusNotFound, nil)
		w.WriteHeader(errRes.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUser)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	if err := h.service.Delete(userID); err != nil {
		errRes := api.NewErrorResponse(err.Error(), http.StatusNotFound, nil)
		w.WriteHeader(errRes.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errRes)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
