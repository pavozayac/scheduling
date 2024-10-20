package model

import (
	"errors"
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

type Input struct {
	scheduleId int
	firstArg   int
	secondArg  int
	thirdArg   ConstraintType
}

type Output struct {
	expectedConstraint Constraint
	expectedError      error
}

func TestShouldConstructValidPairConstraints(t *testing.T) {
	t.Run("NewLocationTaskConstraint", func(t *testing.T) {
		constraint, err := NewLocationTaskConstraint(123, 2345, 3456, Must)

		assert.Equal(t, Constraint{123, -1, 3456, 2345, -1, -1, Must}, constraint)
		assert.Nil(t, err)
	})

	t.Run("NewTaskWorkerConstraint", func(t *testing.T) {
		constraint, err := NewTaskWorkerConstraint(123, 2345, 3456, Must)

		assert.Equal(t, Constraint{123, 2345, 3456, -1, -1, -1, Must}, constraint)
		assert.Nil(t, err)
	})

	t.Run("NewLocationWorkerConstraint", func(t *testing.T) {
		constraint, err := NewLocationWorkerConstraint(123, 2345, 3456, Must)

		assert.Equal(t, Constraint{123, 3456, -1, 2345, -1, -1, Must}, constraint)
		assert.Nil(t, err)
	})
}

func TestShouldThrowOnInvalidPairConstraints(t *testing.T) {
	var testcases = []struct {
		input  Input
		output Output
	}{
		{
			Input{-1, 2000, 3000, Must},
			Output{Constraint{}, shared.ErrNegativeId},
		},
		{
			Input{1234, -1, 3000, Cannot},
			Output{Constraint{}, shared.ErrNegativeId},
		},
		{
			Input{1234, 141, -1, Must},
			Output{Constraint{}, shared.ErrNegativeId},
		},
	}

	var constructors = []struct {
		function func(int, int, int, ConstraintType) (Constraint, error)
		name     string
	}{
		{NewTaskWorkerConstraint, "NewTaskWorkerConstraint"},
		{NewLocationTaskConstraint, "NewLocationTaskConstraint"},
		{NewLocationWorkerConstraint, "NewLocationWorkerConstraint"},
	}

	for _, testcase := range testcases {
		for _, constructor := range constructors {
			t.Run(constructor.name, func(t *testing.T) {
				c, err := constructor.function(testcase.input.scheduleId, testcase.input.firstArg, testcase.input.secondArg, testcase.input.thirdArg)

				assert.Equal(t, testcase.output.expectedConstraint, c)
				assert.ErrorIs(t, err, testcase.output.expectedError)
			})
		}
	}
}

func TestShouldThrowOnInvalidTimeConstraints(t *testing.T) {
	type TimeInput struct {
		scheduleId int
		id         int
		startTime  int
		endTime    int
	}

	testcases := []struct {
		input  TimeInput
		output Output
	}{
		{
			TimeInput{scheduleId: -1, id: 1, startTime: 1, endTime: 2},
			Output{Constraint{}, errors.New("invalid arguments")},
		},
		{
			TimeInput{scheduleId: 1, id: -1, startTime: 1, endTime: 2},
			Output{Constraint{}, errors.New("invalid arguments")},
		},
		{
			TimeInput{scheduleId: 1, id: 1, startTime: 2, endTime: 1},
			Output{Constraint{}, errors.New("invalid arguments")},
		},
	}

	constructors := []struct {
		function func(int, int, int, int, ConstraintType) (Constraint, error)
		name     string
	}{
		{NewWorkerTimeConstraint, "NewWorkerTimeConstraint"},
		{NewTaskTimeConstraint, "NewTaskTimeConstraint"},
		{NewLocationTimeConstraint, "NewLocationTimeConstraint"},
	}

	for _, testcase := range testcases {
		for _, constructor := range constructors {
			t.Run(constructor.name, func(t *testing.T) {
				c, err := constructor.function(testcase.input.scheduleId, testcase.input.id, testcase.input.startTime, testcase.input.endTime, Must)

				assert.Equal(t, testcase.output.expectedConstraint, c)
				assert.ErrorIs(t, err, shared.ErrInvalidArguments)
			})
		}
	}
}

func TestShouldConstructValidTimeConstraints(t *testing.T) {
	t.Run("NewWorkerTimeConstraint", func(t *testing.T) {
		constraint, err := NewWorkerTimeConstraint(123, 2345, 3456, 4567, Must)

		assert.Equal(t, Constraint{123, 2345, -1, -1, 3456, 4567, Must}, constraint)
		assert.Nil(t, err)
	})

	t.Run("NewTaskTimeConstraint", func(t *testing.T) {
		constraint, err := NewTaskTimeConstraint(123, 2345, 3456, 4567, Must)

		assert.Equal(t, Constraint{123, -1, 2345, -1, 3456, 4567, Must}, constraint)
		assert.Nil(t, err)
	})

	t.Run("NewLocationTimeConstraint", func(t *testing.T) {
		constraint, err := NewLocationTimeConstraint(123, 2345, 3456, 4567, Must)

		assert.Equal(t, Constraint{123, -1, -1, 2345, 3456, 4567, Must}, constraint)
		assert.Nil(t, err)
	})
}

func TestShouldDetectConflicts(t *testing.T) {
	testcases := []struct {
		name     string
		c1       Constraint
		c2       Constraint
		expected bool
	}{
		{
			name:     "No Conflict - Different ScheduleId",
			c1:       Constraint{123, 2345, -1, -1, 3456, 4567, Must},
			c2:       Constraint{124, 2345, -1, -1, 3456, 4567, Must},
			expected: false,
		},
		{
			name:     "Conflict - Overlapping Time",
			c1:       Constraint{123, 2345, -1, -1, 3456, 4567, Must},
			c2:       Constraint{123, 2345, -1, -1, 4000, 5000, Must},
			expected: true,
		},
		{
			name:     "No Conflict - Non-overlapping Time",
			c1:       Constraint{123, 2345, -1, -1, 3456, 4567, Must},
			c2:       Constraint{123, 2345, -1, -1, 4568, 5678, Must},
			expected: false,
		},
		{
			name:     "Conflict - Duplicate Constraint",
			c1:       Constraint{123, 2345, -1, -1, 3456, 4567, Must},
			c2:       Constraint{123, 2345, -1, -1, 3456, 4567, Must},
			expected: true,
		},
		{
			name:     "Conflict - Mutually Exclusive Constraints",
			c1:       Constraint{123, 2345, -1, -1, 3456, 4567, Must},
			c2:       Constraint{123, 2345, -1, -1, 3456, 4567, Cannot},
			expected: true,
		},
		{
			name:     "No Conflict - Different WorkerId",
			c1:       Constraint{123, 2345, -1, -1, 3456, 4567, Must},
			c2:       Constraint{123, 2346, -1, -1, 3456, 4567, Must},
			expected: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			result := testcase.c1.ConflictsWith(testcase.c2)
			assert.Equal(t, testcase.expected, result)
		})
	}
}
