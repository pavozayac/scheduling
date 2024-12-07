package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
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
		id.Scan(dto.Id)
	}

	scheduleId.Scan(dto.ScheduleId)

	worker, err := model.NewWorker(id, scheduleId, dto.FirstName, dto.LastName)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveOrUpdateWorker(ctx, *worker)
	if err != nil {
		return nil, err
	}

	feedback, err := s.GetWorker(ctx, dto.Id)
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

func (s *workerService) GetWorker(ctx context.Context, id string) (*ports.WorkerDTO, error) {
	var workerId shared.Identity
	workerId.Scan(id)

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
	workerId.Scan(id)

	return s.repo.DeleteWorker(ctx, workerId)
}

var _ ports.WorkerService = &workerService{}
