package services

import (
	"context"
	"testing"

	"github.com/pavozayac/constraints/internal/domain/model"
	"github.com/pavozayac/constraints/internal/domain/ports"
	"github.com/pavozayac/constraints/internal/domain/shared"
	"github.com/pavozayac/constraints/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTaskService(t *testing.T) {
	repo := mocks.NewMockTaskRepository(t)
	service := NewTaskService(repo)

	uuidOne := shared.MockIdentityGenerator{}.Generate()
	uuidTwo := shared.MockIdentityGenerator{}.Generate()

	exampleTask, err := model.NewTask(uuidOne, uuidTwo, "Task A", "Description A")
	require.NoError(t, err)

	t.Run("ShouldCreateNewTask", func(t *testing.T) {
		saveCall := repo.On("SaveOrUpdateTask", mock.Anything, mock.Anything).Return(nil)
		getCall := repo.On("GetTask", mock.Anything, mock.Anything).Return(exampleTask, nil).NotBefore(saveCall)
		t.Cleanup(func() {
			saveCall.Unset()
			getCall.Unset()
		})

		task, err := service.CreateOrModifyTask(context.Background(), ports.TaskDTO{
			Id:          uuidOne.String(),
			ScheduleId:  uuidTwo.String(),
			Title:       "Task A",
			Description: "Description A",
		})

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), task.Id)
		assert.Equal(t, uuidTwo.String(), task.ScheduleId)
		assert.Equal(t, "Task A", task.Title)
		assert.Equal(t, "Description A", task.Description)
	})

	t.Run("ShouldGetTask", func(t *testing.T) {
		getCall := repo.On("GetTask", mock.Anything, uuidOne).Return(exampleTask, nil)
		t.Cleanup(func() {
			getCall.Unset()
		})

		task, err := service.GetTask(context.Background(), uuidOne.String())

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), task.Id)
		assert.Equal(t, uuidTwo.String(), task.ScheduleId)
		assert.Equal(t, "Task A", task.Title)
		assert.Equal(t, "Description A", task.Description)
	})

	t.Run("ShouldRemoveTask", func(t *testing.T) {
		deleteCall := repo.On("DeleteTask", mock.Anything, uuidOne).Return(nil)
		t.Cleanup(func() {
			deleteCall.Unset()
		})

		err := service.RemoveTask(context.Background(), uuidOne.String())

		repo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("ShouldUpdateTask", func(t *testing.T) {
		updatedTask, err := model.NewTask(uuidOne, uuidTwo, "Task B", "Description B")
		require.NoError(t, err)

		saveCall := repo.On("SaveOrUpdateTask", mock.Anything, mock.Anything).Return(nil)
		getCall := repo.On("GetTask", mock.Anything, mock.Anything).Return(updatedTask, nil).NotBefore(saveCall)
		t.Cleanup(func() {
			saveCall.Unset()
			getCall.Unset()
		})

		task, err := service.CreateOrModifyTask(context.Background(), ports.TaskDTO{
			Id:          uuidOne.String(),
			ScheduleId:  uuidTwo.String(),
			Title:       "Task B",
			Description: "Description B",
		})

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), task.Id)
		assert.Equal(t, uuidTwo.String(), task.ScheduleId)
		assert.Equal(t, "Task B", task.Title)
		assert.Equal(t, "Description B", task.Description)
	})
}
