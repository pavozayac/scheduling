package ports

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type TaskRepository interface {
	SaveOrUpdateTask(model.Task) error

	GetTask(shared.Identity) (model.Task, error)

	RemoveTask(shared.Identity) error
}
