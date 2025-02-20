package types

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ResultResponse struct {
	Result string `json:"result"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type PostTaskResponse struct {
	TaskId string `json:"task_id"`
}

type PostTaskRequestBody struct {
	Code         string `json:"code"`
	CompilerName string `json:"compiler_name"`
}

func ExtractPostTaskBody(r *http.Request) (*PostTaskRequestBody, error) {
	data := &PostTaskRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, errors.New("request data invalid")
	}

	return data, ValidatePostTaskBody(data)
}

func ValidatePostTaskBody(data *PostTaskRequestBody) error {
	if data.Code == "" {
		return errors.New("empty code")
	}

	if data.CompilerName == "" {
		return errors.New("empty compiler name")
	}

	return nil
}
