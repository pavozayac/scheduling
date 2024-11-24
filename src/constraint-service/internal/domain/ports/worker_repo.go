package ports

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type WorkerRepository interface {
	SaveOrUpdateWorker(model.Worker) error

	GetWorker(shared.Identity) (model.Worker, error)

	RemoveWorker(shared.Identity) error
}
