package shared

import "errors"

var ErrNilIdentity = errors.New("id arguments must not be nil")
var ErrInvalidArguments = errors.New("invalid arguments")
var ErrConflictingConstraint = errors.New("constraint conflicts with existing constraint")
var ErrNotFound = errors.New("object not found")
