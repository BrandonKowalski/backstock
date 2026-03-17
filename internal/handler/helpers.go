package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError logs the real error and returns a generic message to the client.
func writeError(w http.ResponseWriter, msg string, status int, err error) {
	log.Printf("ERROR: %s: %v", msg, err)
	http.Error(w, msg, status)
}
