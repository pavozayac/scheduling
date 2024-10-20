package model

import (
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateNewSchedule(t *testing.T) {
	testcases := []struct {
		name             string
		id               int
		title            string
		constraints      Constraints
		expectedError    error
		expectedSchedule *Schedule
	}{
		{
			"Valid Schedule",
			1,
			"Schedule 1",
			Constraints{},
			nil,
			&Schedule{1, "Schedule 1", Constraints{}},
		},
		{
			"Invalid Schedule - Negative ID",
			-1,
			"Schedule 2",
			Constraints{},
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Schedule - Empty Title",
			2,
			"",
			Constraints{},
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Schedule - Nil Constraints",
			3,
			"Schedule 3",
			nil,
			shared.ErrInvalidArguments,
			nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			schedule, err := NewSchedule(testcase.id, testcase.title, testcase.constraints)

			assert.Equal(t, testcase.expectedSchedule, schedule)
			assert.ErrorIs(t, err, testcase.expectedError)
		})
	}
}

func TestShouldIndicateScheduleEquality(t *testing.T) {
	testcases := []struct {
		name     string
		s1       *Schedule
		s2       *Schedule
		expected bool
	}{
		{
			"Equal Schedules",
			&Schedule{1, "Schedule 1", Constraints{}},
			&Schedule{1, "Schedule 1", Constraints{}},
			true,
		},
		{
			"Different IDs",
			&Schedule{1, "Schedule 1", Constraints{}},
			&Schedule{2, "Schedule 2", Constraints{}},
			false,
		},
		{
			"Different IDs, same data",
			&Schedule{1, "Schedule 1", Constraints{}},
			&Schedule{2, "Schedule 1", Constraints{}},
			false,
		},
		{
			"Nil Schedule",
			nil,
			&Schedule{1, "Schedule 1", Constraints{}},
			false,
		},
		{
			"Both Nil Schedules",
			nil,
			nil,
			false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			result := testcase.s1.Equals(testcase.s2)
			assert.Equal(t, testcase.expected, result)
		})
	}
}

func TestShouldAddConstraint(t *testing.T) {
	schedule := &Schedule{1, "Schedule 1", Constraints{}}
	constraint := Constraint{ /*...*/ }

	err := schedule.AddConstraint(constraint)
	assert.Nil(t, err)
	assert.Contains(t, schedule.constraints, constraint)
}

func TestShouldNotAddDuplicateConstraint(t *testing.T) {
	constraint := Constraint{ /*...*/ }
	schedule := &Schedule{1, "Schedule 1", Constraints{constraint}}

	err := schedule.AddConstraint(constraint)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrConflictingConstraint)
}

func TestShouldNotAddConflictingConstraint(t *testing.T) {
	constraint1 := Constraint{ /*...*/ }
	constraint2 := Constraint{ /*...*/ }
	schedule := &Schedule{1, "Schedule 1", Constraints{constraint1}}

	// Assuming ConflictsWith is implemented
	err := schedule.AddConstraint(constraint2)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrConflictingConstraint)
}

func TestShouldRemoveConstraint(t *testing.T) {
	constraint := Constraint{ /*...*/ }
	schedule := &Schedule{1, "Schedule 1", Constraints{constraint}}

	err := schedule.RemoveConstraint(constraint)
	assert.Nil(t, err)
	assert.NotContains(t, schedule.constraints, constraint)
}

func TestShouldNotRemoveNonExistentConstraint(t *testing.T) {
	constraint := Constraint{ /*...*/ }
	schedule := &Schedule{1, "Schedule 1", Constraints{}}

	err := schedule.RemoveConstraint(constraint)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrNotFound)
}
