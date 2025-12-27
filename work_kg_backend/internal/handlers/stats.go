package handlers

import (
	"encoding/json"
	"net/http"

	"work_kg_backend/internal/database"
)

func HandleGetStats(w http.ResponseWriter, r *http.Request) {
	stats := database.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
