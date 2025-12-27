package handlers

import (
	"encoding/json"
	"net/http"

	"work_kg_backend/internal/database"
	"work_kg_backend/internal/models"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := database.GetAdminByEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !database.ValidateAdminPassword(user.Password, req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Return email as token (simple auth - use JWT in production)
	response := models.LoginResponse{
		Token: user.Email,
		User: map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleGetMe(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("X-User-Email")

	user, err := database.GetAdminByEmailWithoutPassword(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
