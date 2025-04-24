package usecases

import "codeproccesor/src/domain"

type Processor interface {
	Process(domain.Task) error
}
