package types

import (
	"encoding/json"
	"errors"
	"net/http"
)

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func ExtractRegisterBody(r *http.Request) (*RegisterRequestBody, error) {
	data := &RegisterRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, errors.New("request data invalid")
	}

	return data, nil
}

func ExtractLoginBody(r *http.Request) (*LoginRequestBody, error) {
	data := &LoginRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, errors.New("request data invalid")
	}

	return data, nil
}
