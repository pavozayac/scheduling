package model

import (
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

type Input struct {
	scheduleId shared.Identity
	firstArg   shared.Identity
	secondArg  shared.Identity
	thirdArg   ConstraintType
}

type Output struct {
	expectedConstraint Constraint
	expectedError      error
}

var mockId1 = shared.MockIdentityGenerator{}.Generate()
var mockId2 = shared.MockIdentityGenerator{}.Generate()
var mockId3 = shared.MockIdentityGenerator{}.Generate()
var mockId4 = shared.MockIdentityGenerator{}.Generate()

func TestShouldConstructValidPairConstraints(t *testing.T) {
	t.Run("NewLocationTaskConstraint", func(t *testing.T) {
		constraint, err := NewLocationTaskConstraint(mockId1, mockId2, mockId3, Must)

		assert.Equal(t, Constraint{mockId1, shared.NilIdentity, mockId3, mockId2, -1, -1, Must}, constraint)
		assert.Nil(t, err)
	})

	t.Run("NewTaskWorkerConstraint", func(t *testing.T) {
		constraint, err := NewTaskWorkerConstraint(mockId1, mockId2, mockId3, Must)

		assert.Equal(t, Constraint{mockId1, mockId2, mockId3, shared.NilIdentity, -1, -1, Must}, constraint)
		assert.Nil(t, err)
	})

	t.Run("NewLocationWorkerConstraint", func(t *testing.T) {
		constraint, err := NewLocationWorkerConstraint(mockId1, mockId2, mockId3, Must)

		assert.Equal(t, Constraint{mockId1, mockId3, shared.NilIdentity, mockId2, -1, -1, Must}, constraint)
		assert.Nil(t, err)
	})
}

func TestShouldThrowOnInvalidPairConstraints(t *testing.T) {
	var testcases = []struct {
		input  Input
		output Output
	}{
		{
			Input{shared.NilIdentity, mockId2, mockId3, Must},
			Output{Constraint{}, shared.ErrNilIdentity},
		},
		{
			Input{mockId1, shared.NilIdentity, mockId3, Cannot},
			Output{Constraint{}, shared.ErrNilIdentity},
		},
		{
			Input{mockId1, mockId2, shared.NilIdentity, Must},
			Output{Constraint{}, shared.ErrNilIdentity},
		},
	}

	var constructors = []struct {
		function func(shared.Identity, shared.Identity, shared.Identity, ConstraintType) (Constraint, error)
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
		scheduleId shared.Identity
		id         shared.Identity
		startTime  int
		endTime    int
	}

	testcases := []struct {
		input  TimeInput
		output Output
	}{
		{
			TimeInput{mockId1, mockId2, -1, 2},
			Output{Constraint{}, shared.ErrInvalidArguments},
		},
		{
			TimeInput{mockId1, mockId2, 1, -1},
			Output{Constraint{}, shared.ErrInvalidArguments},
		},
		{
			TimeInput{mockId1, mockId2, 2, -1},
			Output{Constraint{}, shared.ErrInvalidArguments},
		},
	}

	constructors := []struct {
		function func(shared.Identity, shared.Identity, int, int, ConstraintType) (Constraint, error)
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
				assert.ErrorIs(t, err, testcase.output.expectedError)
			})
		}
	}
}

func TestShouldConstructValidTimeConstraints(t *testing.T) {
	t.Run("NewWorkerTimeConstraint", func(t *testing.T) {
		constraint, err := NewWorkerTimeConstraint(mockId1, mockId2, 3456, 4567, Must)

		assert.Equal(t, Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must}, constraint)
		assert.Nil(t, err)
	})

	t.Run("NewTaskTimeConstraint", func(t *testing.T) {
		constraint, err := NewTaskTimeConstraint(mockId1, mockId2, 3456, 4567, Must)

		assert.Equal(t, Constraint{mockId1, shared.NilIdentity, mockId2, shared.NilIdentity, 3456, 4567, Must}, constraint)
		assert.Nil(t, err)
	})

	t.Run("NewLocationTimeConstraint", func(t *testing.T) {
		constraint, err := NewLocationTimeConstraint(mockId1, mockId2, 3456, 4567, Must)

		assert.Equal(t, Constraint{mockId1, shared.NilIdentity, shared.NilIdentity, mockId2, 3456, 4567, Must}, constraint)
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
			c1:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must},
			c2:       Constraint{mockId3, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must},
			expected: false,
		},
		{
			name:     "Conflict - Overlapping Time",
			c1:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must},
			c2:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 4000, 5000, Must},
			expected: true,
		},
		{
			name:     "No Conflict - Non-overlapping Time",
			c1:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must},
			c2:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 4568, 5678, Must},
			expected: false,
		},
		{
			name:     "Conflict - Duplicate Constraint",
			c1:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must},
			c2:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must},
			expected: true,
		},
		{
			name:     "Conflict - Mutually Exclusive Constraints",
			c1:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must},
			c2:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Cannot},
			expected: true,
		},
		{
			name:     "No Conflict - Different WorkerId",
			c1:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must},
			c2:       Constraint{mockId1, mockId3, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must},
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
