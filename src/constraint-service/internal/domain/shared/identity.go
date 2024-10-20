package shared

import "github.com/google/uuid"

type Identity uuid.UUID

var NilIdentity = Identity(uuid.Nil)

type IdentityGenerator interface {
	Generate() Identity
}

type UuidGenerator struct{}

func (g UuidGenerator) Generate() Identity {
	return Identity(uuid.New())
}

func (g UuidGenerator) GetNull() Identity {
	return Identity(uuid.Nil)
}

type MockIdentityGenerator struct {
}

func (g MockIdentityGenerator) Generate() Identity {
	return Identity(uuid.New())
}
