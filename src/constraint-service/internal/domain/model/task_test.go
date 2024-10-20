package model

import (
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateNewTask(t *testing.T) {
	testcases := []struct {
		name          string
		id            int
		taskName      string
		description   string
		scheduleID    int
		expectedError error
		expectedTask  *Task
	}{
		{
			"Valid Task",
			1,
			"Task A",
			"Description A",
			1,
			nil,
			&Task{1, "Task A", "Description A", 1},
		},
		{
			"Invalid Task - Negative ID",
			-1,
			"Task B",
			"Description B",
			2,
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Task - Empty Name",
			2,
			"",
			"Description C",
			3,
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Task - Empty Description",
			3,
			"Task D",
			"",
			4,
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Task - Negative Schedule ID",
			4,
			"Task E",
			"Description E",
			-1,
			shared.ErrInvalidArguments,
			nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			task, err := NewTask(testcase.id, testcase.scheduleID, testcase.taskName, testcase.description)

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
			&Task{1, "Task A", "Description A", 1},
			&Task{1, "Task A", "Description A", 1},
			true,
		},
		{
			"Different IDs",
			&Task{1, "Task A", "Description A", 1},
			&Task{2, "Task B", "Description B", 2},
			false,
		},
		{
			"Different IDs, same data",
			&Task{1, "Task A", "Description A", 1},
			&Task{2, "Task A", "Description A", 1},
			false,
		},
		{
			"Nil Task",
			nil,
			&Task{1, "Task A", "Description A", 1},
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
