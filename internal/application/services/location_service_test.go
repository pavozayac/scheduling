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

func TestLocationService(t *testing.T) {
	repo := mocks.NewMockLocationRepository(t)
	service := NewLocationService(repo)

	uuidOne := shared.MockIdentityGenerator{}.Generate()
	uuidTwo := shared.MockIdentityGenerator{}.Generate()

	exampleLocation, err := model.NewLocation(uuidOne, uuidTwo, "some name", "some description")
	require.NoError(t, err)

	t.Run("ShouldCreateNewLocation", func(t *testing.T) {
		saveCall := repo.On("SaveOrUpdateLocation", mock.Anything, mock.Anything).Return(nil)
		getCall := repo.On("GetLocation", mock.Anything, mock.Anything).Return(exampleLocation, nil).NotBefore(saveCall)
		t.Cleanup(func() {
			saveCall.Unset()
			getCall.Unset()
		})

		location, err := service.CreateOrModifyLocation(context.Background(), ports.LocationDTO{
			Id:          uuidOne.String(),
			ScheduleId:  uuidTwo.String(),
			Name:        "some name",
			Description: "some description",
		})

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), location.Id)
		assert.Equal(t, uuidTwo.String(), location.ScheduleId)
		assert.Equal(t, "some name", location.Name)
		assert.Equal(t, "some description", location.Description)
	})

	t.Run("ShouldGetLocation", func(t *testing.T) {
		getCall := repo.On("GetLocation", mock.Anything, uuidOne).Return(exampleLocation, nil)
		t.Cleanup(func() {
			getCall.Unset()
		})

		location, err := service.GetLocation(context.Background(), uuidOne.String())

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), location.Id)
		assert.Equal(t, uuidTwo.String(), location.ScheduleId)
		assert.Equal(t, "some name", location.Name)
		assert.Equal(t, "some description", location.Description)
	})

	t.Run("ShouldRemoveLocation", func(t *testing.T) {
		deleteCall := repo.On("DeleteLocation", mock.Anything, uuidOne).Return(nil)
		t.Cleanup(func() {
			deleteCall.Unset()
		})

		err := service.RemoveLocation(context.Background(), uuidOne.String())

		repo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("ShouldUpdateLocation", func(t *testing.T) {
		updatedLocation, err := model.NewLocation(uuidOne, uuidTwo, "updated name", "updated description")
		require.NoError(t, err)

		saveCall := repo.On("SaveOrUpdateLocation", mock.Anything, mock.Anything).Return(nil)
		getCall := repo.On("GetLocation", mock.Anything, mock.Anything).Return(updatedLocation, nil).NotBefore(saveCall)
		t.Cleanup(func() {
			saveCall.Unset()
			getCall.Unset()
		})

		location, err := service.CreateOrModifyLocation(context.Background(), ports.LocationDTO{
			Id:          uuidOne.String(),
			ScheduleId:  uuidTwo.String(),
			Name:        "updated name",
			Description: "updated description",
		})

		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, uuidOne.String(), location.Id)
		assert.Equal(t, uuidTwo.String(), location.ScheduleId)
		assert.Equal(t, "updated name", location.Name)
		assert.Equal(t, "updated description", location.Description)
	})
}
