package model

import (
	"github.com/pavozayac/constraints/internal/domain/shared"
)

type Location struct {
	id          shared.Identity
	scheduleId  shared.Identity
	name        string
	description string
}

func (l Location) Id() shared.Identity {
	return l.id
}

func (l Location) ScheduleId() shared.Identity {
	return l.scheduleId
}

func (l Location) Name() string {
	return l.name
}

func (l Location) Description() string {
	return l.description
}

func NewLocation(id, scheduleId shared.Identity, name string, description string) (*Location, error) {
	if id == shared.NilIdentity || name == "" || description == "" || scheduleId == shared.NilIdentity {
		return nil, shared.ErrInvalidArguments
	}

	return &Location{
		id:          id,
		scheduleId:  scheduleId,
		name:        name,
		description: description,
	}, nil
}

func (l *Location) Equals(other *Location) bool {
	if l == nil || other == nil {
		return false
	}
	return l.id == other.id
}

type Locations []*Location
