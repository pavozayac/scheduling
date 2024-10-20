package model

import (
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

func TestNewLocation(t *testing.T) {
	testcases := []struct {
		name             string
		id               int
		locationName     string
		description      string
		scheduleId       int
		expectedError    error
		expectedLocation *Location
	}{
		{
			name:             "Valid Location",
			id:               1,
			locationName:     "Location A",
			description:      "Description A",
			scheduleId:       100,
			expectedError:    nil,
			expectedLocation: &Location{id: 1, name: "Location A", description: "Description A", scheduleId: 100},
		},
		{
			name:             "Invalid Location - Negative ID",
			id:               -1,
			locationName:     "Location B",
			description:      "Description B",
			scheduleId:       101,
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
		{
			name:             "Invalid Location - Empty Name",
			id:               2,
			locationName:     "",
			description:      "Description C",
			scheduleId:       102,
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
		{
			name:             "Invalid Location - Empty Description",
			id:               3,
			locationName:     "Location D",
			description:      "",
			scheduleId:       103,
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
		{
			name:             "Invalid Location - Negative ScheduleId",
			id:               4,
			locationName:     "Location E",
			description:      "Description E",
			scheduleId:       -104,
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			location, err := NewLocation(testcase.id, testcase.locationName, testcase.description, testcase.scheduleId)

			assert.Equal(t, testcase.expectedLocation, location)
			assert.ErrorIs(t, err, testcase.expectedError)
		})
	}
}

func TestLocation_Equals(t *testing.T) {
	testcases := []struct {
		name     string
		l1       *Location
		l2       *Location
		expected bool
	}{
		{
			name:     "Equal Locations",
			l1:       &Location{id: 1, name: "Location A", description: "Description A", scheduleId: 100},
			l2:       &Location{id: 1, name: "Location A", description: "Description A", scheduleId: 100},
			expected: true,
		},
		{
			name:     "Different IDs",
			l1:       &Location{id: 1, name: "Location A", description: "Description A", scheduleId: 100},
			l2:       &Location{id: 2, name: "Location B", description: "Description B", scheduleId: 101},
			expected: false,
		},
		{
			name:     "Different IDs, same data",
			l1:       &Location{id: 1, name: "Location A", description: "Description A", scheduleId: 100},
			l2:       &Location{id: 2, name: "Location A", description: "Description A", scheduleId: 100},
			expected: false,
		},
		{
			name:     "Nil Location",
			l1:       nil,
			l2:       &Location{id: 1, name: "Location A", description: "Description A", scheduleId: 100},
			expected: false,
		},
		{
			name:     "Both Nil Locations",
			l1:       nil,
			l2:       nil,
			expected: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			result := testcase.l1.Equals(testcase.l2)
			assert.Equal(t, testcase.expected, result)
		})
	}
}
