package processpool

import "errors"

var (
	ErrProcessPoolStopped = errors.New("processpool was stopped")
)
