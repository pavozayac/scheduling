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

func TestWorkerService(t *testing.T) {
	repo := mocks.NewMockWorkerRepository(t)
	service := NewWorkerService(repo)

	uuidOne := shared.MockIdentityGenerator{}.Generate()
	uuidTwo := shared.MockIdentityGenerator{}.Generate()

	exampleWorker, err := model.NewWorker(uuidOne, uuidTwo, "John", "Doe")
	require.NoError(t, err)

	t.Run("ShouldCreateNewWorker", func(t *testing.T) {
		saveCall := repo.On("SaveOrUpdateWorker", mock.Anything, mock.Anything).Return(nil)
		getCall := repo.On("GetWorker", mock.Anything, mock.Anything).Return(exampleWorker, nil).NotBefore(saveCall)
		t.Cleanup(func() {
			saveCall.Unset()
			getCall.Unset()
		})

		worker, err := service.CreateOrModifyWorker(context.Background(), ports.WorkerDTO{
			Id:         uuidOne.String(),
			ScheduleId: uuidTwo.String(),
			FirstName:  "John",
			LastName:   "Doe",
		})

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), worker.Id)
		assert.Equal(t, uuidTwo.String(), worker.ScheduleId)
		assert.Equal(t, "John", worker.FirstName)
		assert.Equal(t, "Doe", worker.LastName)
	})

	t.Run("ShouldGetWorker", func(t *testing.T) {
		getCall := repo.On("GetWorker", mock.Anything, uuidOne).Return(exampleWorker, nil)
		t.Cleanup(func() {
			getCall.Unset()
		})

		worker, err := service.GetWorker(context.Background(), uuidOne.String())

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), worker.Id)
		assert.Equal(t, uuidTwo.String(), worker.ScheduleId)
		assert.Equal(t, "John", worker.FirstName)
		assert.Equal(t, "Doe", worker.LastName)
	})

	t.Run("ShouldRemoveWorker", func(t *testing.T) {
		deleteCall := repo.On("DeleteWorker", mock.Anything, uuidOne).Return(nil)
		t.Cleanup(func() {
			deleteCall.Unset()
		})

		err := service.RemoveWorker(context.Background(), uuidOne.String())

		assert.NoError(t, err)
	})

	t.Run("ShouldUpdateWorker", func(t *testing.T) {
		updatedWorker, err := model.NewWorker(uuidOne, uuidTwo, "Jane", "Doe")
		require.NoError(t, err)

		saveCall := repo.On("SaveOrUpdateWorker", mock.Anything, mock.Anything).Return(nil)
		getCall := repo.On("GetWorker", mock.Anything, mock.Anything).Return(updatedWorker, nil).NotBefore(saveCall)
		t.Cleanup(func() {
			saveCall.Unset()
			getCall.Unset()
		})

		worker, err := service.CreateOrModifyWorker(context.Background(), ports.WorkerDTO{
			Id:         uuidOne.String(),
			ScheduleId: uuidTwo.String(),
			FirstName:  "Jane",
			LastName:   "Doe",
		})
		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), worker.Id)
		assert.Equal(t, uuidTwo.String(), worker.ScheduleId)
		assert.Equal(t, "Jane", worker.FirstName)
		assert.Equal(t, "Doe", worker.LastName)
	})
}
