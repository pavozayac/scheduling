package model

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type Location struct {
	id          shared.Identity
	name        string
	description string
	scheduleId  shared.Identity
}

func NewLocation(id shared.Identity, name string, description string, scheduleId shared.Identity) (*Location, error) {
	if id == shared.NilIdentity || name == "" || description == "" || scheduleId == shared.NilIdentity {
		return nil, shared.ErrInvalidArguments
	}

	return &Location{
		id:          id,
		name:        name,
		description: description,
		scheduleId:  scheduleId,
	}, nil
}

func (l *Location) Equals(other *Location) bool {
	if l == nil || other == nil {
		return false
	}
	return l.id == other.id
}

type Locations []*Location
