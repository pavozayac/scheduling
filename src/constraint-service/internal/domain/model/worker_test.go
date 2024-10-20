package model

import (
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateNewWorker(t *testing.T) {
	testcases := []struct {
		name           string
		id             int
		firstName      string
		lastName       string
		scheduleID     int
		expectedError  error
		expectedWorker *Worker
	}{
		{
			"Valid Worker",
			1,
			"John",
			"Doe",
			1,
			nil,
			&Worker{1, "John", "Doe", 1},
		},
		{
			"Invalid Worker - Negative ID",
			-1,
			"Jane",
			"Doe",
			2,
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Worker - Empty First Name",
			2,
			"",
			"Doe",
			3,
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Worker - Empty Last Name",
			3,
			"John",
			"",
			4,
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Worker - Negative Schedule ID",
			4,
			"Jane",
			"Doe",
			-1,
			shared.ErrInvalidArguments,
			nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			worker, err := NewWorker(testcase.id, testcase.firstName, testcase.lastName, testcase.scheduleID)

			assert.Equal(t, testcase.expectedWorker, worker)
			assert.ErrorIs(t, err, testcase.expectedError)
		})
	}
}

func TestShouldIndicateWorkerEquality(t *testing.T) {
	testcases := []struct {
		name     string
		w1       *Worker
		w2       *Worker
		expected bool
	}{
		{
			"Equal Workers",
			&Worker{1, "John", "Doe", 1},
			&Worker{1, "John", "Doe", 1},
			true,
		},
		{
			"Different IDs",
			&Worker{1, "John", "Doe", 1},
			&Worker{2, "Jane", "Doe", 2},
			false,
		},
		{
			"Different IDs, same data",
			&Worker{1, "John", "Doe", 1},
			&Worker{2, "John", "Doe", 1},
			false,
		},
		{
			"Nil Worker",
			nil,
			&Worker{1, "John", "Doe", 1},
			false,
		},
		{
			"Both Nil Workers",
			nil,
			nil,
			false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			result := testcase.w1.Equals(testcase.w2)
			assert.Equal(t, testcase.expected, result)
		})
	}
}
