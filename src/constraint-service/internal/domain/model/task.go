package model

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type Task struct {
	id          int
	name        string
	description string
	scheduleId  int
}

func NewTask(id, scheduleId int, name, description string) (*Task, error) {
	if id < 0 || scheduleId < 0 || name == "" || description == "" {
		return nil, shared.ErrInvalidArguments
	}

	return &Task{
		id:          id,
		name:        name,
		description: description,
		scheduleId:  scheduleId,
	}, nil
}

func (t *Task) Equals(other *Task) bool {
	if t == nil || other == nil {
		return false
	}

	return t.id == other.id
}

type Tasks []*Task
