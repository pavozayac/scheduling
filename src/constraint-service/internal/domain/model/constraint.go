package model

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
)

type ConstraintType string

const (
	Must   ConstraintType = "must"
	Cannot ConstraintType = "cannot"
)

type Constraint struct {
	scheduleId     shared.Identity
	workerId       shared.Identity
	taskId         shared.Identity
	locationId     shared.Identity
	startTime      int
	endTime        int
	constraintType ConstraintType
}

func newConstraint(scheduleId, workerId, taskId, locationId shared.Identity, startTime, endTime int, constraintType ConstraintType) Constraint {
	return Constraint{
		scheduleId:     scheduleId,
		workerId:       workerId,
		taskId:         taskId,
		locationId:     locationId,
		startTime:      startTime,
		endTime:        endTime,
		constraintType: constraintType,
	}
}

func NewTaskWorkerConstraint(scheduleId, workerId, taskId shared.Identity, constraintType ConstraintType) (Constraint, error) {
	if scheduleId == shared.NilIdentity || workerId == shared.NilIdentity || taskId == shared.NilIdentity {
		return Constraint{}, shared.ErrNegativeId
	}
	return newConstraint(scheduleId, workerId, taskId, shared.NilIdentity, -1, -1, constraintType), nil
}

func NewLocationTaskConstraint(scheduleId, locationId, taskId shared.Identity, constraintType ConstraintType) (Constraint, error) {
	if scheduleId == shared.NilIdentity || locationId == shared.NilIdentity || taskId == shared.NilIdentity {
		return Constraint{}, shared.ErrNegativeId
	}
	return newConstraint(scheduleId, shared.NilIdentity, taskId, locationId, -1, -1, constraintType), nil
}

func NewLocationWorkerConstraint(scheduleId, locationId, workerId shared.Identity, constraintType ConstraintType) (Constraint, error) {
	if scheduleId == shared.NilIdentity || locationId == shared.NilIdentity || workerId == shared.NilIdentity {
		return Constraint{}, shared.ErrNegativeId
	}
	return newConstraint(scheduleId, workerId, shared.NilIdentity, locationId, -1, -1, constraintType), nil
}

func NewLocationTimeConstraint(scheduleId, locationId shared.Identity, startTime, endTime int, constraintType ConstraintType) (Constraint, error) {
	if scheduleId == shared.NilIdentity || locationId == shared.NilIdentity {
		return Constraint{}, shared.ErrNegativeId
	}
	if startTime >= endTime || startTime < 0 || endTime < 0 {
		return Constraint{}, shared.ErrInvalidArguments
	}
	return newConstraint(scheduleId, shared.NilIdentity, shared.NilIdentity, locationId, startTime, endTime, constraintType), nil
}

func NewWorkerTimeConstraint(scheduleId, workerId shared.Identity, startTime, endTime int, constraintType ConstraintType) (Constraint, error) {
	if scheduleId == shared.NilIdentity || workerId == shared.NilIdentity {
		return Constraint{}, shared.ErrNegativeId
	}
	if startTime >= endTime || startTime < 0 || endTime < 0 {
		return Constraint{}, shared.ErrInvalidArguments
	}
	return newConstraint(scheduleId, workerId, shared.NilIdentity, shared.NilIdentity, startTime, endTime, constraintType), nil
}

func NewTaskTimeConstraint(scheduleId, taskId shared.Identity, startTime, endTime int, constraintType ConstraintType) (Constraint, error) {
	if scheduleId == shared.NilIdentity || taskId == shared.NilIdentity {
		return Constraint{}, shared.ErrNegativeId
	}
	if startTime >= endTime || startTime < 0 || endTime < 0 {
		return Constraint{}, shared.ErrInvalidArguments
	}
	return newConstraint(scheduleId, shared.NilIdentity, taskId, shared.NilIdentity, startTime, endTime, constraintType), nil
}

func (c Constraint) ConflictsWith(other Constraint) bool {
	if c.Equals(other) {
		return true
	}

	if c.scheduleId != other.scheduleId {
		return false
	}

	if c.workerId != shared.NilIdentity && c.workerId == other.workerId && c.taskId != shared.NilIdentity && c.taskId == other.taskId {
		return true
	}

	if c.locationId != shared.NilIdentity && c.locationId == other.workerId && c.taskId != shared.NilIdentity && c.taskId == other.taskId {
		return true
	}

	if c.workerId != shared.NilIdentity && c.workerId == other.workerId && c.locationId != shared.NilIdentity && c.locationId == other.locationId {
		return true
	}

	if (c.locationId != shared.NilIdentity && c.locationId == other.locationId || c.taskId != shared.NilIdentity && c.taskId == other.taskId || c.workerId != shared.NilIdentity && c.workerId == other.workerId) &&
		c.startTime != -1 && other.startTime != -1 && c.startTime <= other.endTime && c.endTime >= other.startTime {
		return true
	}

	return false
}

func (c Constraint) Equals(other Constraint) bool {
	return c.scheduleId == other.scheduleId &&
		c.workerId == other.workerId &&
		c.taskId == other.taskId &&
		c.locationId == other.locationId &&
		c.startTime == other.startTime &&
		c.endTime == other.endTime &&
		c.constraintType == other.constraintType
}

type Constraints []Constraint
