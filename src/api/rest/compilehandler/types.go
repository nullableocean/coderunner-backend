package compilehandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"nullableocean-postupashki/src/api/rest"
	"nullableocean-postupashki/src/storage"
)

type ResultResponse struct {
	Result string `json:"result"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type CreateTaskResponse struct {
	TaskId string `json:"task_id"`
}

type CreateTaskRequestBody struct {
	Code         string `json:"code"`
	CompilerName string `json:"compiler_name"`
}

func extractCreateTaskBody(r *http.Request) (*CreateTaskRequestBody, error) {
	data := &CreateTaskRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, errors.New("request data invalid")
	}

	return data, validateCreateTaskBody(data)
}

func validateCreateTaskBody(data *CreateTaskRequestBody) error {
	if data.Code == "" {
		return errors.New("empty code")
	}

	if data.CompilerName == "" {
		return errors.New("empty compiler name")
	}

	return nil
}

func proccessError(w http.ResponseWriter, err error, defaultStatus int) {
	status := defaultStatus

	if errors.Is(err, storage.ErrNotFound) {
		status = http.StatusNotFound
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(rest.ErrorResponse{Error: err.Error()})
}
