package model

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type Worker struct {
	id         shared.Identity
	scheduleId shared.Identity
	firstName  string
	lastName   string
}

func NewWorker(id shared.Identity, scheduleId shared.Identity, firstName, lastName string) (*Worker, error) {
	if id == shared.NilIdentity || scheduleId == shared.NilIdentity || firstName == "" || lastName == "" {
		return nil, shared.ErrInvalidArguments
	}

	return &Worker{
		id:         id,
		scheduleId: scheduleId,
		firstName:  firstName,
		lastName:   lastName,
	}, nil
}

func (w *Worker) Equals(other *Worker) bool {
	if w == nil || other == nil {
		return false
	}

	return w.id == other.id
}

type Workers []*Worker
