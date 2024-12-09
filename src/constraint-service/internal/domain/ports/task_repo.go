package ports

import (
	"context"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type TaskRepository interface {
	SaveOrUpdateTask(context.Context, model.Task) error

	GetTask(context.Context, shared.Identity) (*model.Task, error)

	DeleteTask(context.Context, shared.Identity) error
}
