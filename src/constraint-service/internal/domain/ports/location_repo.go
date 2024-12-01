package ports

import (
	"context"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type LocationRepository interface {
	SaveOrUpdateLocation(context.Context, model.Location) error

	GetLocation(context.Context, shared.Identity) (model.Location, error)

	RemoveLocation(context.Context, shared.Identity) error
}
