package ports

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type ScheduleRepository interface {
	SaveOrUpdateSchedule(model.Schedule) error

	GetSchedule(shared.Identity) (model.Schedule, error)

	RemoveSchedule(shared.Identity) error
}
