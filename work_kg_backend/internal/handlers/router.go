package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/login", HandleLogin).Methods("POST")
	api.HandleFunc("/auth/me", AuthMiddleware(HandleGetMe)).Methods("GET")

	// Jobs routes
	api.HandleFunc("/jobs", HandleGetJobs).Methods("GET")
	api.HandleFunc("/jobs", AuthMiddleware(HandleCreateJob)).Methods("POST")
	api.HandleFunc("/jobs/{id}", AuthMiddleware(HandleUpdateJob)).Methods("PUT")
	api.HandleFunc("/jobs/{id}", AuthMiddleware(HandleDeleteJob)).Methods("DELETE")

	// Users routes
	api.HandleFunc("/users", AuthMiddleware(HandleGetUsers)).Methods("GET")

	// Resumes routes
	api.HandleFunc("/resumes", AuthMiddleware(HandleGetResumes)).Methods("GET")

	// Stats route
	api.HandleFunc("/stats", AuthMiddleware(HandleGetStats)).Methods("GET")

	return r
}

func StartServer(port string) {
	r := SetupRouter()

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://work_kg.okugula.dev", "http://localhost:7040", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Printf("HTTP server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
