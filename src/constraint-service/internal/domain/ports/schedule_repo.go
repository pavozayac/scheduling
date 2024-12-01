package ports

import (
	"context"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type ScheduleRepository interface {
	SaveOrUpdateSchedule(context.Context, model.Schedule) error

	GetSchedule(context.Context, shared.Identity) (model.Schedule, error)

	RemoveSchedule(context.Context, shared.Identity) error
}
