package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type scheduleService struct {
	repo ports.ScheduleRepository
}

func NewScheduleService(repo ports.ScheduleRepository) *scheduleService {
	return &scheduleService{repo: repo}
}

func (s *scheduleService) CreateOrModifySchedule(ctx context.Context, dto ports.ScheduleDTO) (*ports.ScheduleDTO, error) {
	var scheduleId, workerId, taskId, locationId shared.Identity
	var constraints []model.Constraint

	if dto.Id == "" {
		scheduleId = shared.UuidGenerator{}.Generate()
	} else {
		if err := scheduleId.Scan(dto.Id); err != nil {
			return nil, err
		}
	}

	for _, c := range dto.Constraints {
		err1 := workerId.Scan(c.WorkerId)
		err2 := taskId.Scan(c.TaskId)
		err3 := locationId.Scan(c.LocationId)

		if err1 != nil || err2 != nil || err3 != nil {
			return nil, errors.Join(err1, err2, err3)
		}

		constraint, err := model.NewConstraint(
			scheduleId,
			workerId,
			taskId,
			locationId,
			c.StartTime,
			c.EndTime,
			model.ConstraintType(c.ConstraintType),
		)
		if err != nil {
			return nil, fmt.Errorf("model.NewConstraint() error: %w", err)
		}

		constraints = append(constraints, constraint)
	}

	schedule, err := model.NewSchedule(scheduleId, dto.Title, constraints)
	if err != nil {
		return nil, fmt.Errorf("model.NewSchedule() error: %w", err)
	}

	err = s.repo.SaveOrUpdateSchedule(ctx, *schedule)
	if err != nil {
		return nil, fmt.Errorf("scheduleService.ScheduleRepository.SaveOrUpdateSchedule() error: %w", err)
	}

	feedback, err := s.GetSchedule(ctx, scheduleId.String())
	if err != nil {
		return nil, fmt.Errorf("scheduleService.ScheduleRepository.GetSchedule() error: %w", err)
	}

	return feedback, nil
}

func (s *scheduleService) GetSchedule(ctx context.Context, id string) (*ports.ScheduleDTO, error) {
	var scheduleId shared.Identity
	if err := scheduleId.Scan(id); err != nil {
		return nil, fmt.Errorf("scheduleId.Scan() error: %w", err)
	}

	schedule, err := s.repo.GetSchedule(ctx, scheduleId)
	if err != nil {
		return nil, fmt.Errorf("scheduleService.ScheduleRepository.GetSchedule() error: %w", err)
	}

	dtoConstraints := make([]ports.ConstraintDTO, 0, len(schedule.Constraints()))
	for _, c := range schedule.Constraints() {
		dtoConstraints = append(dtoConstraints, ports.ConstraintDTO{
			WorkerId:       uuid.UUID(c.WorkerId()).String(),
			TaskId:         uuid.UUID(c.TaskId()).String(),
			LocationId:     uuid.UUID(c.LocationId()).String(),
			StartTime:      c.StartTime(),
			EndTime:        c.EndTime(),
			ConstraintType: string(c.Type()),
		})
	}

	return &ports.ScheduleDTO{
		Id:          uuid.UUID(schedule.Id()).String(),
		Title:       schedule.Title(),
		Constraints: dtoConstraints,
	}, nil
}

func (s *scheduleService) RemoveSchedule(ctx context.Context, id string) error {
	var scheduleId shared.Identity
	if err := scheduleId.Scan(id); err != nil {
		return fmt.Errorf("scheduleId.Scan() error: %w", err)
	}

	return s.repo.DeleteSchedule(ctx, scheduleId)
}

var _ ports.ScheduleService = &scheduleService{}
