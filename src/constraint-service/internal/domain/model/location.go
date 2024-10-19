package model

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type Location struct {
	id          int
	name        string
	description string
	scheduleId  int
}

func NewLocation(id int, name string, description string, scheduleId int) (*Location, error) {
	if id <= 0 || name == "" || description == "" || scheduleId <= 0 {
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
