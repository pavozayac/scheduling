package model

import "errors"

type Schedule struct {
	id          int
	title       string
	constraints Constraints
}

func NewSchedule(id int, title string, constraints Constraints) (*Schedule, error) {
	if id < 0 || title == "" || constraints == nil {
		return nil, errors.New("id, title, locations, tasks, workers, and constraints must not be nullish")
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
		if c.Equals(constraint) {
			return errors.New("constraint already exists")
		}

		if c.ConflictsWith(constraint) {
			return errors.New("constraint conflicts with existing constraint")
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

	return errors.New("constraint does not exist")
}
