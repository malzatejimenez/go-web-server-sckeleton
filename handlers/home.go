package handlers

import (
	"encoding/json"
	"net/http"
	"platzi/go/rest-ws/server"
)

type HomeResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type header
		w.Header().Set("Content-Type", "application/json")

		// Set status code 200 (the request was successful)
		w.WriteHeader(http.StatusOK)

		// Create new response
		response := &HomeResponse{
			Message: "Hello world",
			Status:  true,
		}

		// Encode response
		json.NewEncoder(w).Encode(response)
	}
}
