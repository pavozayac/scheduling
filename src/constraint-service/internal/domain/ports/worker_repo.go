package ports

import (
	"context"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type WorkerRepository interface {
	SaveOrUpdateWorker(context.Context, model.Worker) error

	GetWorker(context.Context, shared.Identity) (*model.Worker, error)

	DeleteWorker(context.Context, shared.Identity) error
}
