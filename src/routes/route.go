package routes

import (
	"github.com/gorilla/mux"
)

func ApiV1Routes(r *mux.Router) {
	v1 := r.PathPrefix("/api/v1").Subrouter()
	HealthRoutes(v1)
}
