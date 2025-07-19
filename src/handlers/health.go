package handlers

import (
	"encoding/json"
	"net/http"
)

type BasicHealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := BasicHealthResponse{
			Status:  "healthy",
			Message: "Service is running",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
