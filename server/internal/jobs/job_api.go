package jobs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/agierlicki/loom/internal/api"
	"github.com/google/uuid"
)

type JobApi struct {
	store JobStore
}

func NewJobApi(store JobStore) *JobApi {
	return &JobApi{
		store: store,
	}
}

func (a *JobApi) JobsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.GetJobsHandler(w, r)
	case "POST":
		a.CreateJobHandler(w, r)
	default:
		api.JSONError(w, "Method not implemented", http.StatusNotImplemented)
	}
}

func (a *JobApi) JobHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.GetJobHandler(w, r)
	default:
		api.JSONError(w, "Method not implemented", http.StatusNotImplemented)
	}
}

func (a *JobApi) GetJobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := a.store.List()

	if err != nil {
		api.JSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	api.JSONResponse(w, jobs)
}

func (a *JobApi) GetJobHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := a.store.Get(id)

	if err != nil {
		api.JSONError(w, fmt.Sprintf("Job with id %q not found!", id), http.StatusNotFound)
		return
	}

	api.JSONResponse(w, job)
}

func (a *JobApi) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	var req JobCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	newJob := Job{
		ID:          uuid.NewString(),
		Type:        req.Type,
		Payload:     req.Payload,
		Status:      JobPending,
		Attempts:    0,
		MaxAttempts: 3, // erstmal hart codiert, später ggf. konfigurierbar
		CreatedAt:   time.Now(),
	}

	if err := a.store.Create(newJob); err != nil {
		api.JSONError(w, "could not create job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	api.JSONResponse(w, newJob)
}
