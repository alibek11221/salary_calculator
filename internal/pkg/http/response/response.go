package response

import (
	"encoding/json"
	"net/http"
)

func Ok(w http.ResponseWriter, data any) error {
	return WriteJSON(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data any) error {
	return WriteJSON(w, http.StatusCreated, data)
}

func BadRequest(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusBadRequest, message)
}

func Conflict(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusConflict, message)
}

func InternalServerError(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusInternalServerError, message)
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, map[string]string{"error": message})
}
