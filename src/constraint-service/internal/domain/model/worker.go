package model

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type Worker struct {
	id         int
	firstName  string
	lastName   string
	scheduleId int
}

func NewWorker(id int, firstName, lastName string, scheduleId int) (*Worker, error) {
	if id < 0 || scheduleId < 0 || firstName == "" || lastName == "" {
		return nil, shared.ErrInvalidArguments
	}

	return &Worker{
		id:         id,
		firstName:  firstName,
		lastName:   lastName,
		scheduleId: scheduleId,
	}, nil
}

func (w *Worker) Equals(other *Worker) bool {
	if w == nil || other == nil {
		return false
	}

	return w.id == other.id
}

type Workers []*Worker
