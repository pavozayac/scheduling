package model

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type Task struct {
	id          shared.Identity
	scheduleId  shared.Identity
	name        string
	description string
}

func (t Task) Id() shared.Identity {
	return t.id
}

func (t Task) ScheduleId() shared.Identity {
	return t.scheduleId
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Description() string {
	return t.description
}

func NewTask(id, scheduleId shared.Identity, name, description string) (*Task, error) {
	if id == shared.NilIdentity || scheduleId == shared.NilIdentity || name == "" || description == "" {
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
