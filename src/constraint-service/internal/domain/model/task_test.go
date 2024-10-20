package model

import (
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateNewTask(t *testing.T) {
	testcases := []struct {
		name          string
		id            shared.Identity
		scheduleId    shared.Identity
		taskName      string
		description   string
		expectedError error
		expectedTask  *Task
	}{
		{
			"Valid Task",
			mockId1,
			mockId2,
			"Task A",
			"Description A",
			nil,
			&Task{mockId1, mockId2, "Task A", "Description A"},
		},
		{
			"Invalid Task - Nil ID",
			shared.NilIdentity,
			mockId1,
			"Task B",
			"Description B",
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Task - Empty Name",
			mockId1,
			mockId2,
			"",
			"Description C",
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Task - Empty Description",
			mockId1,
			mockId2,
			"Task D",
			"",
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Task - Negative Schedule ID",
			mockId1,
			shared.NilIdentity,
			"Task E",
			"Description E",
			shared.ErrInvalidArguments,
			nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			task, err := NewTask(testcase.id, testcase.scheduleId, testcase.taskName, testcase.description)

			assert.Equal(t, testcase.expectedTask, task)
			assert.ErrorIs(t, err, testcase.expectedError)
		})
	}
}

func TestShouldIndicateEquality(t *testing.T) {
	testcases := []struct {
		name     string
		t1       *Task
		t2       *Task
		expected bool
	}{
		{
			"Equal Tasks",
			&Task{mockId1, mockId1, "Task A", "Description A"},
			&Task{mockId1, mockId1, "Task A", "Description A"},
			true,
		},
		{
			"Different IDs",
			&Task{mockId1, mockId1, "Task A", "Description A"},
			&Task{mockId2, mockId2, "Task B", "Description B"},
			false,
		},
		{
			"Different IDs, same data",
			&Task{mockId1, mockId1, "Task A", "Description A"},
			&Task{mockId2, mockId1, "Task A", "Description A"},
			false,
		},
		{
			"Nil Task",
			nil,
			&Task{mockId1, mockId1, "Task A", "Description A"},
			false,
		},
		{
			"Both Nil Tasks",
			nil,
			nil,
			false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			result := testcase.t1.Equals(testcase.t2)
			assert.Equal(t, testcase.expected, result)
		})
	}
}
