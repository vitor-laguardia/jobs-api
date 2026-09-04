package job

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vitor-laguardia/jobs-api/internal/shared/httputil"
)

type JobHandler struct {
	service *JobService
	router  *http.ServeMux
}

func NewJobHandler(service *JobService) *JobHandler {
	jh := &JobHandler{service: service}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jh.postJob)
	mux.HandleFunc("PUT /jobs/{id}", jh.updateJob)
	mux.HandleFunc("GET /jobs/{id}", jh.getJob)
	mux.HandleFunc("DELETE /jobs/{id}", jh.deleteJob)
	jh.router = mux
	return jh
}

func (jh *JobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jh.router.ServeHTTP(w, r)
}

func (jh *JobHandler) getJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")

	job, err := jh.service.GetByID(jobID)

	if err != nil {
		errRes := httputil.NewErrorResponse(err.Error(), http.StatusNotFound, nil)
		w.WriteHeader(errRes.Status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (jh *JobHandler) postJob(w http.ResponseWriter, r *http.Request) {
	jobReq, errResp := httputil.DecodeValid[CreateJobRequest](w, r)

	if errResp != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errResp.Status)
		json.NewEncoder(w).Encode(errResp)
		return
	}

	nj := jh.service.Create(jobReq)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nj)
}

func (jh *JobHandler) updateJob(w http.ResponseWriter, r *http.Request) {
	jobReq, errResp := httputil.DecodeValid[UpdateJobRequest](w, r)

	if errResp != nil {
		w.WriteHeader(errResp.Status)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	jobID := r.PathValue("id")
	updatedJob, err := jh.service.Update(jobID, jobReq)

	if err != nil {
		var errRes *httputil.ErrorResponse
		switch {
		case errors.Is(err, ErrJobNotFound):
			errRes = httputil.NewErrorResponse(err.Error(), http.StatusBadRequest, nil)
		case errors.Is(err, ErrJobStatusTransition):
			errRes = httputil.NewErrorResponse(err.Error(), http.StatusUnprocessableEntity, nil)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(errRes.Status)
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(updatedJob)
}

func (jh *JobHandler) deleteJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")

	if err := jh.service.Delete(jobID); err != nil {
		errRes := httputil.NewErrorResponse(err.Error(), http.StatusNotFound, nil)
		w.WriteHeader(errRes.Status)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
