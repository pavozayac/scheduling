package shared

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type Identity uuid.UUID

func (i Identity) Value() (driver.Value, error) {
	return uuid.UUID(i).Value()
}

func (i Identity) String() string {
	return uuid.UUID(i).String()
}

func (i Identity) Bytes() ([]byte, error) {
	return uuid.UUID(i).MarshalBinary()
}

func (i *Identity) Scan(src interface{}) error {
	return (*uuid.UUID)(i).Scan(src)
}

var NilIdentity = Identity(uuid.Nil)

type IdentityGenerator interface {
	Generate() Identity
	FromBytes([]byte) (Identity, error)
}

type UuidGenerator struct{}

func (g UuidGenerator) Generate() Identity {
	return Identity(uuid.New())
}

func (g UuidGenerator) FromBytes(b []byte) (Identity, error) {
	id, err := uuid.FromBytes(b)
	if err != nil {
		return NilIdentity, err
	}

	return Identity(id), nil
}

type MockIdentityGenerator struct {
}

func (g MockIdentityGenerator) Generate() Identity {
	return Identity(uuid.New())
}
