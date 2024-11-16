package ports

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type LocationRepository interface {
	SaveOrUpdateLocation(model.Location) error

	GetLocation(shared.Identity) (model.Location, error)

	RemoveLocation(shared.Identity) error
}
