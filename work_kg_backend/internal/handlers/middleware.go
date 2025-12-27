package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"work_kg_backend/internal/database"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Simple token validation - in production use JWT
		token = strings.TrimPrefix(token, "Bearer ")

		user, err := database.GetAdminByEmailWithoutPassword(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", strconv.FormatInt(user.ID, 10))
		r.Header.Set("X-User-Email", user.Email)
		r.Header.Set("X-User-Role", user.Role)

		next(w, r)
	}
}
