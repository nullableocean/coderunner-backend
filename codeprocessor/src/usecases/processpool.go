package usecases

import "codeproccesor/src/domain"

type ProcessPool interface {
	Push(domain.Task) error
	IsRunned() bool
	Stop()
}
