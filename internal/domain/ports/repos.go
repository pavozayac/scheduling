package ports

import (
	"context"

	"github.com/pavozayac/constraints/internal/domain/model"
	"github.com/pavozayac/constraints/internal/domain/shared"
)

type ScheduleRepository interface {
	SaveOrUpdateSchedule(context.Context, model.Schedule) error

	GetSchedule(context.Context, shared.Identity) (*model.Schedule, error)

	DeleteSchedule(context.Context, shared.Identity) error
}

type LocationRepository interface {
	SaveOrUpdateLocation(context.Context, model.Location) error

	GetLocation(context.Context, shared.Identity) (*model.Location, error)

	DeleteLocation(context.Context, shared.Identity) error
}

type TaskRepository interface {
	SaveOrUpdateTask(context.Context, model.Task) error

	GetTask(context.Context, shared.Identity) (*model.Task, error)

	DeleteTask(context.Context, shared.Identity) error
}

type WorkerRepository interface {
	SaveOrUpdateWorker(context.Context, model.Worker) error

	GetWorker(context.Context, shared.Identity) (*model.Worker, error)

	DeleteWorker(context.Context, shared.Identity) error
}
