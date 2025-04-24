package types

import (
	"encoding/json"
	"errors"
	"net/http"
)

type PostCommitRequestBody struct {
	TaskUuid  string `json:"task_uuid"`
	Output    string `json:"output"`
	Error     string `json:"error"`
	IsSuccess bool   `json:"is_success"`
}

func ExtractPostResultBody(r *http.Request) (*PostCommitRequestBody, error) {
	data := &PostCommitRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, errors.New("request data invalid")
	}

	return data, nil
}
