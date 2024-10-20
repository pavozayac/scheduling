package model

import (
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/stretchr/testify/assert"
)

func TestNewLocation(t *testing.T) {
	testcases := []struct {
		name             string
		id               shared.Identity
		locationName     string
		description      string
		scheduleId       shared.Identity
		expectedError    error
		expectedLocation *Location
	}{
		{
			name:             "Valid Location",
			id:               mockId1,
			locationName:     "Location A",
			description:      "Description A",
			scheduleId:       mockId2,
			expectedError:    nil,
			expectedLocation: &Location{id: mockId1, name: "Location A", description: "Description A", scheduleId: mockId2},
		},
		{
			name:             "Invalid Location - Negative ID",
			id:               shared.NilIdentity,
			locationName:     "Location B",
			description:      "Description B",
			scheduleId:       mockId2,
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
		{
			name:             "Invalid Location - Empty Name",
			id:               mockId1,
			locationName:     "",
			description:      "Description C",
			scheduleId:       mockId2,
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
		{
			name:             "Invalid Location - Empty Description",
			id:               mockId1,
			locationName:     "Location D",
			description:      "",
			scheduleId:       mockId2,
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
		{
			name:             "Invalid Location - Negative ScheduleId",
			id:               mockId1,
			locationName:     "Location E",
			description:      "Description E",
			scheduleId:       shared.NilIdentity,
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
			l1:       &Location{id: mockId1, name: "Location A", description: "Description A", scheduleId: mockId2},
			l2:       &Location{id: mockId1, name: "Location A", description: "Description A", scheduleId: mockId2},
			expected: true,
		},
		{
			name:     "Different IDs",
			l1:       &Location{id: mockId1, name: "Location A", description: "Description A", scheduleId: mockId2},
			l2:       &Location{id: mockId3, name: "Location B", description: "Description B", scheduleId: mockId4},
			expected: false,
		},
		{
			name:     "Different IDs, same data",
			l1:       &Location{id: mockId1, name: "Location A", description: "Description A", scheduleId: mockId2},
			l2:       &Location{id: mockId3, name: "Location A", description: "Description A", scheduleId: mockId2},
			expected: false,
		},
		{
			name:     "Nil Location",
			l1:       nil,
			l2:       &Location{id: mockId1, name: "Location A", description: "Description A", scheduleId: mockId2},
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
