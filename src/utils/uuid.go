package utils

import "github.com/google/uuid"

func CreateUuid() string {
	return uuid.NewString()
}
