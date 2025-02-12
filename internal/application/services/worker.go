package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/pavozayac/constraints/internal/domain/model"
	"github.com/pavozayac/constraints/internal/domain/ports"
	"github.com/pavozayac/constraints/internal/domain/shared"
)

type workerService struct {
	repo ports.WorkerRepository
}

func NewWorkerService(repo ports.WorkerRepository) *workerService {
	return &workerService{repo: repo}
}

func (s *workerService) CreateOrModifyWorker(ctx context.Context, dto ports.WorkerDTO) (*ports.WorkerDTO, error) {
	var id, scheduleId shared.Identity

	if dto.Id == "" {
		id = shared.UuidGenerator{}.Generate()
	} else {
		if err := id.Scan(dto.Id); err != nil {
			return nil, err
		}
	}

	if err := scheduleId.Scan(dto.ScheduleId); err != nil {
		return nil, err
	}

	worker, err := model.NewWorker(id, scheduleId, dto.FirstName, dto.LastName)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveOrUpdateWorker(ctx, *worker)
	if err != nil {
		return nil, err
	}

	feedback, err := s.GetWorker(ctx, id.String())
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

func (s *workerService) GetWorker(ctx context.Context, id string) (*ports.WorkerDTO, error) {
	var workerId shared.Identity
	if err := workerId.Scan(id); err != nil {
		return nil, err
	}

	worker, err := s.repo.GetWorker(ctx, workerId)
	if err != nil {
		return nil, err
	}

	return &ports.WorkerDTO{
		Id:         uuid.UUID(worker.Id()).String(),
		ScheduleId: uuid.UUID(worker.ScheduleId()).String(),
		FirstName:  worker.FirstName(),
		LastName:   worker.LastName(),
	}, nil
}

func (s *workerService) RemoveWorker(ctx context.Context, id string) error {
	var workerId shared.Identity
	if err := workerId.Scan(id); err != nil {
		return err
	}

	return s.repo.DeleteWorker(ctx, workerId)
}

var _ ports.WorkerService = &workerService{}
