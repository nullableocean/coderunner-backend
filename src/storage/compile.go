package storage

import (
	"nullableocean-postupashki/src/models"
)

type CompileStorage interface {
	SaveTask(task *models.CompileTask) error
	FindTask(id string) (*models.CompileTask, error)
}
