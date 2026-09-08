package job

import (
	"encoding/json"
	"errors"
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
	mux.HandleFunc("POST /jobs", h.postJob)
	mux.HandleFunc("PUT /jobs/{id}", h.updateJob)
	mux.HandleFunc("GET /jobs/{id}", h.getJob)
	mux.HandleFunc("DELETE /jobs/{id}", h.deleteJob)
	h.router = mux
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) getJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")

	job, err := h.service.GetByID(jobID)

	if err != nil {
		errRes := api.NewErrorResponse(err.Error(), http.StatusNotFound, nil)
		w.WriteHeader(errRes.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (h *Handler) postJob(w http.ResponseWriter, r *http.Request) {
	jobReq, errResp := api.DecodeValid[CreateJobRequest](w, r)

	if errResp != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errResp.Status)
		json.NewEncoder(w).Encode(errResp)
		return
	}

	nj := h.service.Create(jobReq)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nj)
}

func (h *Handler) updateJob(w http.ResponseWriter, r *http.Request) {
	jobReq, errResp := api.DecodeValid[UpdateJobRequest](w, r)

	if errResp != nil {
		w.WriteHeader(errResp.Status)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	jobID := r.PathValue("id")
	updatedJob, err := h.service.Update(jobID, jobReq)

	if err != nil {
		var errRes *api.ErrorResponse
		switch {
		case errors.Is(err, ErrJobNotFound):
			errRes = api.NewErrorResponse(err.Error(), http.StatusBadRequest, nil)
		case errors.Is(err, ErrJobStatusTransition):
			errRes = api.NewErrorResponse(err.Error(), http.StatusUnprocessableEntity, nil)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(errRes.Status)
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(updatedJob)
}

func (h *Handler) deleteJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")

	if err := h.service.Delete(jobID); err != nil {
		errRes := api.NewErrorResponse(err.Error(), http.StatusNotFound, nil)
		w.WriteHeader(errRes.Status)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
