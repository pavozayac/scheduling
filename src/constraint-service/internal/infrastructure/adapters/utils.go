package adapters

import (
	"github.com/google/uuid"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

func nillifyIdentity(id shared.Identity) *uuid.UUID {
	if id == shared.NilIdentity {
		return nil
	}
	identity := uuid.UUID(id)
	return &identity
}

func denillifyUuid(uiid *uuid.UUID) shared.Identity {
	if uiid == nil {
		return shared.NilIdentity
	}

	return shared.Identity(*uiid)
}
