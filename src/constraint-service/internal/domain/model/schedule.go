package model

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type Schedule struct {
	id          int
	title       string
	constraints Constraints
}

func NewSchedule(id int, title string, constraints Constraints) (*Schedule, error) {
	if id < 0 || title == "" || constraints == nil {
		return nil, shared.ErrInvalidArguments
	}

	return &Schedule{
		id:          id,
		title:       title,
		constraints: constraints,
	}, nil
}

func (s *Schedule) Equals(other *Schedule) bool {
	if s == nil || other == nil {
		return false
	}

	return s.id == other.id
}

func (s *Schedule) AddConstraint(constraint Constraint) error {
	for _, c := range s.constraints {
		if c.Equals(constraint) || c.ConflictsWith(constraint) {
			return shared.ErrConflictingConstraint
		}
	}

	s.constraints = append(s.constraints, constraint)

	return nil
}

func (s *Schedule) RemoveConstraint(constraint Constraint) error {
	for i, c := range s.constraints {
		if c.Equals(constraint) {
			s.constraints = append(s.constraints[:i], s.constraints[i+1:]...)
			return nil
		}
	}

	return shared.ErrNotFound
}
