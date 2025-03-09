package types

import (
	"encoding/json"
	"errors"
	"net/http"
	"nullableocean-postupashki/src/repository"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ProcessError(w http.ResponseWriter, err error, defaultStatus int) {
	status := defaultStatus

	if errors.Is(err, repository.ErrTaskNotFound) {
		status = http.StatusNotFound
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}
