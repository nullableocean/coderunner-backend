package repository

import "nullableocean-postupashki/src/domain"

type TaskSender interface {
	Send(*domain.Task) error
}
