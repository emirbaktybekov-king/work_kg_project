package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"work_kg_backend/internal/database"
	"work_kg_backend/internal/models"
)

func HandleGetJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := database.GetAllJobs()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	job.Source = "admin"
	job.IsActive = true

	if err := database.CreateJob(&job); err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

func HandleUpdateJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := database.UpdateJob(id, &job); err != nil {
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		return
	}

	job.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func HandleDeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	if err := database.DeleteJob(id); err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
