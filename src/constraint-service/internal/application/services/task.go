package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type taskService struct {
	repo ports.TaskRepository
}

func NewTaskService(repo ports.TaskRepository) *taskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateOrModifyTask(ctx context.Context, dto ports.TaskDTO) (*ports.TaskDTO, error) {
	var id, scheduleId shared.Identity

	if dto.Id == "" {
		id = shared.UuidGenerator{}.Generate()
	} else {
		id.Scan(dto.Id)
	}

	scheduleId.Scan(dto.ScheduleId)

	task, err := model.NewTask(id, scheduleId, dto.Title, dto.Description)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveOrUpdateTask(ctx, *task)
	if err != nil {
		return nil, err
	}

	feedback, err := s.GetTask(ctx, id.String())
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

func (s *taskService) GetTask(ctx context.Context, id string) (*ports.TaskDTO, error) {
	var taskId shared.Identity
	taskId.Scan(id)

	task, err := s.repo.GetTask(ctx, taskId)
	if err != nil {
		return nil, err
	}

	return &ports.TaskDTO{
		Id:          uuid.UUID(task.Id()).String(),
		ScheduleId:  uuid.UUID(task.ScheduleId()).String(),
		Title:       task.Name(),
		Description: task.Description(),
	}, nil
}

func (s *taskService) RemoveTask(ctx context.Context, id string) error {
	var taskId shared.Identity
	taskId.Scan(id)

	return s.repo.DeleteTask(ctx, taskId)
}

var _ ports.TaskService = &taskService{}
