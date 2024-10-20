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
		scheduleId       shared.Identity
		locationName     string
		description      string
		expectedError    error
		expectedLocation *Location
	}{
		{
			name:             "Valid Location",
			id:               mockId1,
			scheduleId:       mockId2,
			locationName:     "Location A",
			description:      "Description A",
			expectedError:    nil,
			expectedLocation: &Location{id: mockId1, name: "Location A", description: "Description A", scheduleId: mockId2},
		},
		{
			name:             "Invalid Location - Negative ID",
			id:               shared.NilIdentity,
			scheduleId:       mockId2,
			locationName:     "Location B",
			description:      "Description B",
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
		{
			name:             "Invalid Location - Empty Name",
			id:               mockId1,
			scheduleId:       mockId2,
			locationName:     "",
			description:      "Description C",
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
		{
			name:             "Invalid Location - Empty Description",
			id:               mockId1,
			scheduleId:       mockId2,
			locationName:     "Location D",
			description:      "",
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
		{
			name:             "Invalid Location - Negative ScheduleId",
			id:               mockId1,
			scheduleId:       shared.NilIdentity,
			locationName:     "Location E",
			description:      "Description E",
			expectedError:    shared.ErrInvalidArguments,
			expectedLocation: nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			location, err := NewLocation(testcase.id, testcase.scheduleId, testcase.locationName, testcase.description)

			assert.Equal(t, testcase.expectedLocation, location)
			assert.ErrorIs(t, err, testcase.expectedError)
		})
	}
}
