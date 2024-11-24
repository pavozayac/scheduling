package model

import (
	"fmt"

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

func (c Constraint) ScheduleId() shared.Identity {
	return c.scheduleId
}

func (c Constraint) WorkerId() shared.Identity {
	return c.workerId
}

func (c Constraint) TaskId() shared.Identity {
	return c.taskId
}

func (c Constraint) LocationId() shared.Identity {
	return c.locationId
}

func (c Constraint) StartTime() int {
	return c.startTime
}

func (c Constraint) EndTime() int {
	return c.endTime
}

func (c Constraint) Type() ConstraintType {
	return c.constraintType
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

func validateConstraintData(
	scheduleId, workerId, taskId, locationId shared.Identity,
	startTime, endTime int,
	constraintType ConstraintType,
) error {
	if scheduleId == shared.NilIdentity {
		return shared.ErrNilIdentity
	}

	hasWorker := workerId != shared.NilIdentity
	hasTask := taskId != shared.NilIdentity
	hasLocation := locationId != shared.NilIdentity
	hasTime := startTime != -1 && endTime != -1

	// Validate time range if present
	if hasTime {
		if startTime >= endTime || startTime < 0 || endTime < 0 {
			return shared.ErrInvalidArguments
		}
	}

	if constraintType != Must && constraintType != Cannot {
		return shared.ErrInvalidArguments
	}

	// Match valid patterns
	switch {
	case hasWorker && hasTask && !hasLocation && !hasTime: // TaskWorker
		return nil
	case hasLocation && hasTask && !hasWorker && !hasTime: // LocationTask
		return nil
	case hasLocation && hasWorker && !hasTask && !hasTime: // LocationWorker
		return nil
	case hasLocation && hasTime && !hasWorker && !hasTask: // LocationTime
		return nil
	case hasWorker && hasTime && !hasLocation && !hasTask: // WorkerTime
		return nil
	case hasTask && hasTime && !hasLocation && !hasWorker: // TaskTime
		return nil
	default:
		return shared.ErrInvalidArguments
	}
}

func NewConstraint(scheduleId, workerId, taskId, locationId shared.Identity, startTime, endTime int, constraintType ConstraintType) (Constraint, error) {
	if err := validateConstraintData(scheduleId, workerId, taskId, locationId, startTime, endTime, constraintType); err != nil {
		return Constraint{}, fmt.Errorf("constraint validation failed: %w", err)
	}

	return newConstraint(scheduleId, workerId, taskId, locationId, startTime, endTime, constraintType), nil
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
