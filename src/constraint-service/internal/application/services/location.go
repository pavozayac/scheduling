package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type locationService struct {
	repo ports.LocationRepository
}

func NewLocationService(repo ports.LocationRepository) *locationService {
	return &locationService{repo: repo}
}

func (s *locationService) CreateOrModifyLocation(ctx context.Context, dto ports.LocationDTO) (*ports.LocationDTO, error) {
	var id, scheduleId shared.Identity

	if dto.Id == "" {
		id = shared.UuidGenerator{}.Generate()
	} else {
		id.Scan(dto.Id)
	}

	scheduleId.Scan(dto.ScheduleId)

	location, err := model.NewLocation(id, scheduleId, dto.Name, dto.Description)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveOrUpdateLocation(ctx, *location)
	if err != nil {
		return nil, err
	}

	feedback, err := s.GetLocation(ctx, id.String())
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

func (s *locationService) GetLocation(ctx context.Context, id string) (*ports.LocationDTO, error) {
	var locationId shared.Identity
	locationId.Scan(id)

	location, err := s.repo.GetLocation(ctx, locationId)
	if err != nil {
		return nil, err
	}

	result := &ports.LocationDTO{
		Id:          uuid.UUID(location.Id()).String(),
		ScheduleId:  uuid.UUID(location.ScheduleId()).String(),
		Name:        location.Name(),
		Description: location.Description(),
	}

	return result, nil
}

func (s *locationService) RemoveLocation(ctx context.Context, id string) error {
	var locationId shared.Identity
	locationId.Scan(id)

	return s.repo.DeleteLocation(ctx, locationId)
}

var _ ports.LocationService = &locationService{}
