package services

import (
	"context"
	"testing"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestScheduleService(t *testing.T) {
	repo := mocks.NewMockScheduleRepository(t)
	service := NewScheduleService(repo)

	uuidOne := shared.MockIdentityGenerator{}.Generate()
	uuidTwo := shared.MockIdentityGenerator{}.Generate()

	exampleConstraint, err := model.NewConstraint(uuidOne, uuidTwo, shared.NilIdentity, shared.NilIdentity, 1609459200, 1609462800, model.Must)
	require.NoError(t, err)

	exampleSchedule, err := model.NewSchedule(uuidOne, "Test Schedule", []model.Constraint{exampleConstraint})
	require.NoError(t, err)

	t.Run("ShouldCreateNewSchedule", func(t *testing.T) {
		saveCall := repo.On("SaveOrUpdateSchedule", mock.Anything, mock.Anything).Return(nil)
		getCall := repo.On("GetSchedule", mock.Anything, mock.Anything).Return(exampleSchedule, nil).NotBefore(saveCall)
		t.Cleanup(func() {
			saveCall.Unset()
			getCall.Unset()
		})

		schedule, err := service.CreateOrModifySchedule(context.Background(), ports.ScheduleDTO{
			Id:    uuidOne.String(),
			Title: "Test Schedule",
			Constraints: []ports.ConstraintDTO{
				{
					WorkerId:       uuidTwo.String(),
					TaskId:         shared.NilIdentity.String(),
					LocationId:     shared.NilIdentity.String(),
					StartTime:      1609459200,
					EndTime:        1609462800,
					ConstraintType: string(model.Must),
				},
			},
		})

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), schedule.Id)
		assert.Equal(t, "Test Schedule", schedule.Title)
		assert.Equal(t, 1, len(schedule.Constraints))
		assert.Equal(t, uuidTwo.String(), schedule.Constraints[0].WorkerId)
		assert.Equal(t, shared.NilIdentity.String(), schedule.Constraints[0].TaskId)
		assert.Equal(t, shared.NilIdentity.String(), schedule.Constraints[0].LocationId)
		assert.Equal(t, 1609459200, schedule.Constraints[0].StartTime)
		assert.Equal(t, 1609462800, schedule.Constraints[0].EndTime)
		assert.Equal(t, string(model.Must), schedule.Constraints[0].ConstraintType)
	})

	t.Run("ShouldGetSchedule", func(t *testing.T) {
		getCall := repo.On("GetSchedule", mock.Anything, uuidOne).Return(exampleSchedule, nil)
		t.Cleanup(func() {
			getCall.Unset()
		})

		schedule, err := service.GetSchedule(context.Background(), uuidOne.String())

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), schedule.Id)
		assert.Equal(t, "Test Schedule", schedule.Title)
		assert.Equal(t, 1, len(schedule.Constraints))
		assert.Equal(t, uuidTwo.String(), schedule.Constraints[0].WorkerId)
		assert.Equal(t, shared.NilIdentity.String(), schedule.Constraints[0].TaskId)
		assert.Equal(t, shared.NilIdentity.String(), schedule.Constraints[0].LocationId)
		assert.Equal(t, 1609459200, schedule.Constraints[0].StartTime)
		assert.Equal(t, 1609462800, schedule.Constraints[0].EndTime)
		assert.Equal(t, string(model.Must), schedule.Constraints[0].ConstraintType)
	})

	t.Run("ShouldRemoveSchedule", func(t *testing.T) {
		deleteCall := repo.On("DeleteSchedule", mock.Anything, uuidOne).Return(nil)
		t.Cleanup(func() {
			deleteCall.Unset()
		})

		err := service.RemoveSchedule(context.Background(), uuidOne.String())

		repo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("ShouldUpdateSchedule", func(t *testing.T) {
		updatedConstraint, err := model.NewConstraint(uuidOne, uuidTwo, shared.NilIdentity, shared.NilIdentity, 1609459200, 1609462800, model.Must)
		require.NoError(t, err)

		updatedSchedule, err := model.NewSchedule(uuidOne, "Updated Schedule", []model.Constraint{updatedConstraint})
		require.NoError(t, err)

		saveCall := repo.On("SaveOrUpdateSchedule", mock.Anything, mock.Anything).Return(nil)
		getCall := repo.On("GetSchedule", mock.Anything, mock.Anything).Return(updatedSchedule, nil).NotBefore(saveCall)
		t.Cleanup(func() {
			saveCall.Unset()
			getCall.Unset()
		})

		schedule, err := service.CreateOrModifySchedule(context.Background(), ports.ScheduleDTO{
			Id:    uuidOne.String(),
			Title: "Updated Schedule",
			Constraints: []ports.ConstraintDTO{
				{
					WorkerId:       uuidTwo.String(),
					TaskId:         shared.NilIdentity.String(),
					LocationId:     shared.NilIdentity.String(),
					StartTime:      1609459200,
					EndTime:        1609462800,
					ConstraintType: string(model.Must),
				},
			},
		})

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), schedule.Id)
		assert.Equal(t, "Updated Schedule", schedule.Title)
		assert.Equal(t, 1, len(schedule.Constraints))
		assert.Equal(t, uuidTwo.String(), schedule.Constraints[0].WorkerId)
		assert.Equal(t, shared.NilIdentity.String(), schedule.Constraints[0].TaskId)
		assert.Equal(t, shared.NilIdentity.String(), schedule.Constraints[0].LocationId)
		assert.Equal(t, 1609459200, schedule.Constraints[0].StartTime)
		assert.Equal(t, 1609462800, schedule.Constraints[0].EndTime)
		assert.Equal(t, string(model.Must), schedule.Constraints[0].ConstraintType)
	})
}
