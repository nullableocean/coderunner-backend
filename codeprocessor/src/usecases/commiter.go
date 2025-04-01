package usecases

import "codeproccesor/src/domain"

type Commiter interface {
	Commit(domain.Result) error
}
