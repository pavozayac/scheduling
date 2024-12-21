package model

import (
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

var mockId1 = shared.MockIdentityGenerator{}.Generate()
var mockId2 = shared.MockIdentityGenerator{}.Generate()
var mockId3 = shared.MockIdentityGenerator{}.Generate()
var mockId4 = shared.MockIdentityGenerator{}.Generate()

func TestShouldConstructValidConstraints(t *testing.T) {
	testcases := []struct {
		name           string
		scheduleId     shared.Identity
		workerId       shared.Identity
		taskId         shared.Identity
		locationId     shared.Identity
		startTime      int
		endTime        int
		constraintType ConstraintType
		expected       Constraint
	}{
		{
			name:           "LocationTask",
			scheduleId:     mockId1,
			workerId:       shared.NilIdentity,
			taskId:         mockId3,
			locationId:     mockId2,
			startTime:      -1,
			endTime:        -1,
			constraintType: Must,
			expected:       Constraint{mockId1, shared.NilIdentity, mockId3, mockId2, -1, -1, Must},
		},
		{
			name:           "TaskWorker",
			scheduleId:     mockId1,
			workerId:       mockId2,
			taskId:         mockId3,
			locationId:     shared.NilIdentity,
			startTime:      -1,
			endTime:        -1,
			constraintType: Must,
			expected:       Constraint{mockId1, mockId2, mockId3, shared.NilIdentity, -1, -1, Must},
		},
		{
			name:           "LocationWorker",
			scheduleId:     mockId1,
			workerId:       mockId3,
			taskId:         shared.NilIdentity,
			locationId:     mockId2,
			startTime:      -1,
			endTime:        -1,
			constraintType: Must,
			expected:       Constraint{mockId1, mockId3, shared.NilIdentity, mockId2, -1, -1, Must},
		},
		{
			name:           "WorkerTime",
			scheduleId:     mockId1,
			workerId:       mockId2,
			taskId:         shared.NilIdentity,
			locationId:     shared.NilIdentity,
			startTime:      3456,
			endTime:        4567,
			constraintType: Must,
			expected:       Constraint{mockId1, mockId2, shared.NilIdentity, shared.NilIdentity, 3456, 4567, Must},
		},
		{
			name:           "TaskTime",
			scheduleId:     mockId1,
			workerId:       shared.NilIdentity,
			taskId:         mockId2,
			locationId:     shared.NilIdentity,
			startTime:      3456,
			endTime:        4567,
			constraintType: Must,
			expected:       Constraint{mockId1, shared.NilIdentity, mockId2, shared.NilIdentity, 3456, 4567, Must},
		},
		{
			name:           "LocationTime",
			scheduleId:     mockId1,
			workerId:       shared.NilIdentity,
			taskId:         shared.NilIdentity,
			locationId:     mockId2,
			startTime:      3456,
			endTime:        4567,
			constraintType: Must,
			expected:       Constraint{mockId1, shared.NilIdentity, shared.NilIdentity, mockId2, 3456, 4567, Must},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			constraint, err := NewConstraint(
				tc.scheduleId,
				tc.workerId,
				tc.taskId,
				tc.locationId,
				tc.startTime,
				tc.endTime,
				tc.constraintType,
			)

			assert.Equal(t, tc.expected, constraint)
			assert.NoError(t, err)
		})
	}
}

func TestShouldThrowOnInvalidConstraints(t *testing.T) {
	testcases := []struct {
		name           string
		scheduleId     shared.Identity
		workerId       shared.Identity
		taskId         shared.Identity
		locationId     shared.Identity
		startTime      int
		endTime        int
		constraintType ConstraintType
		expectedError  error
	}{
		{
			name:           "Nil ScheduleId",
			scheduleId:     shared.NilIdentity,
			workerId:       mockId2,
			taskId:         mockId3,
			locationId:     shared.NilIdentity,
			startTime:      -1,
			endTime:        -1,
			constraintType: Must,
			expectedError:  shared.ErrNilIdentity,
		},
		{
			name:           "Invalid Time Range",
			scheduleId:     mockId1,
			workerId:       mockId2,
			taskId:         shared.NilIdentity,
			locationId:     shared.NilIdentity,
			startTime:      5,
			endTime:        3,
			constraintType: Must,
			expectedError:  shared.ErrInvalidArguments,
		},
		{
			name:           "Too Many Aspects Defined",
			scheduleId:     mockId1,
			workerId:       mockId2,
			taskId:         mockId3,
			locationId:     mockId4,
			startTime:      -1,
			endTime:        -1,
			constraintType: Must,
			expectedError:  shared.ErrInvalidArguments,
		},
		{
			name:           "Invalid Type",
			scheduleId:     mockId1,
			workerId:       mockId2,
			taskId:         mockId3,
			locationId:     mockId4,
			startTime:      -1,
			endTime:        -1,
			constraintType: "sometype",
			expectedError:  shared.ErrInvalidArguments,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			constraint, err := NewConstraint(
				tc.scheduleId,
				tc.workerId,
				tc.taskId,
				tc.locationId,
				tc.startTime,
				tc.endTime,
				tc.constraintType,
			)

			assert.Equal(t, Constraint{}, constraint)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
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
