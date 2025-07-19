package routes

import (
	"olympus/handlers"

	"github.com/gorilla/mux"
)

func HealthRoutes(r *mux.Router) {
	health := r.PathPrefix("/health").Subrouter()
	health.HandleFunc("/liveness", handlers.HealthHandler()).Methods("GET")
}
