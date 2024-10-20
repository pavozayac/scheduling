package model

import (
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateNewWorker(t *testing.T) {
	testcases := []struct {
		name           string
		id             shared.Identity
		scheduleId     shared.Identity
		firstName      string
		lastName       string
		expectedError  error
		expectedWorker *Worker
	}{
		{
			"Valid Worker",
			mockId1,
			mockId1,
			"John",
			"Doe",
			nil,
			&Worker{mockId1, mockId1, "John", "Doe"},
		},
		{
			"Invalid Worker - Negative ID",
			shared.NilIdentity,
			mockId2,
			"Jane",
			"Doe",
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Worker - Empty First Name",
			mockId2,
			mockId2,
			"",
			"Doe",
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Worker - Empty Last Name",
			mockId2,
			mockId2,
			"John",
			"",
			shared.ErrInvalidArguments,
			nil,
		},
		{
			"Invalid Worker - Negative Schedule ID",
			mockId2,
			shared.NilIdentity,
			"Jane",
			"Doe",
			shared.ErrInvalidArguments,
			nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			worker, err := NewWorker(testcase.id, testcase.scheduleId, testcase.firstName, testcase.lastName)

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
			&Worker{mockId1, mockId1, "John", "Doe"},
			&Worker{mockId1, mockId1, "John", "Doe"},
			true,
		},
		{
			"Different IDs",
			&Worker{mockId1, mockId1, "John", "Doe"},
			&Worker{mockId2, mockId2, "Jane", "Doe"},
			false,
		},
		{
			"Different IDs, same data",
			&Worker{mockId1, mockId1, "John", "Doe"},
			&Worker{mockId2, mockId1, "John", "Doe"},
			false,
		},
		{
			"Nil Worker",
			nil,
			&Worker{mockId1, mockId1, "John", "Doe"},
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
