package handlerutils

import (
	"encoding/json"
	"net/http"
)

type jsonError struct {
	Errors []string `json:"errors"`
}

func newJsonErrors(errors []string) jsonError {
	jsonError := jsonError{
		Errors: errors,
	}

	return jsonError
}

type jsonResponse struct {
	Messages []string `json:"messages"`
}

func newJsonResponse(messages []string) jsonResponse {
	response := jsonResponse{
		Messages: messages,
	}

	return response
}

func ReturnJsonError(w http.ResponseWriter, status int, errors ...string) {
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(newJsonErrors(errors)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func ReturnJsonMessages(w http.ResponseWriter, status int, messages ...string) {
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(newJsonResponse(messages)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func ReturnJson(w http.ResponseWriter, status int, serializable interface{}) {
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(serializable); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
