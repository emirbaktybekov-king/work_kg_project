package handlers

import (
	"encoding/json"
	"net/http"

	"work_kg_backend/internal/database"
)

func HandleGetResumes(w http.ResponseWriter, r *http.Request) {
	resumes, err := database.GetAllResumes()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resumes)
}
