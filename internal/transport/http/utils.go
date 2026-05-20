package http

import (
	"encoding/json"
	"net/http"
)

func readJSON(r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func errorJSON(w http.ResponseWriter, status int, err string) {
	writeJSON(w, status, &ErrorResponse{
		Error: err,
	})
}
