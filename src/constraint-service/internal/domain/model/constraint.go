package model

import (
	"errors"
)

type ConstraintType string

const (
	Must   ConstraintType = "must"
	Cannot ConstraintType = "cannot"
)

type Constraint struct {
	scheduleId     int
	workerId       int
	taskId         int
	locationId     int
	startTime      int
	endTime        int
	constraintType ConstraintType
}

func newConstraint(scheduleId, workerId, taskId, locationId, startTime, endTime int, constraintType ConstraintType) Constraint {
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

var ErrNegativeId = errors.New("id arguments must be greater than or equal to 0")

func NewTaskWorkerConstraint(scheduleId, workerId, taskId int, constraintType ConstraintType) (Constraint, error) {
	if scheduleId < 0 || workerId < 0 || taskId < 0 {
		return Constraint{}, ErrNegativeId
	}
	return newConstraint(scheduleId, workerId, taskId, -1, -1, -1, constraintType), nil
}

func NewLocationTaskConstraint(scheduleId, locationId, taskId int, constraintType ConstraintType) (Constraint, error) {
	if scheduleId < 0 || locationId < 0 || taskId < 0 {
		return Constraint{}, ErrNegativeId
	}
	return newConstraint(scheduleId, -1, taskId, locationId, -1, -1, constraintType), nil
}

func NewLocationWorkerConstraint(scheduleId, locationId, workerId int, constraintType ConstraintType) (Constraint, error) {
	if scheduleId < 0 || locationId < 0 || workerId < 0 {
		return Constraint{}, ErrNegativeId
	}
	return newConstraint(scheduleId, workerId, -1, locationId, -1, -1, constraintType), nil
}

var ErrInvalidArguments = errors.New("scheduleId, workerId, and taskId must not be nullish")

func NewLocationTimeConstraint(scheduleId, locationId, startTime, endTime int, constraintType ConstraintType) (Constraint, error) {
	if scheduleId < 0 || locationId < 0 || startTime >= endTime {
		return Constraint{}, ErrInvalidArguments
	}

	return newConstraint(scheduleId, -1, -1, locationId, startTime, endTime, constraintType), nil
}
func NewWorkerTimeConstraint(scheduleId, workerId, startTime, endTime int, constraintType ConstraintType) (Constraint, error) {
	if scheduleId < 0 || workerId < 0 || startTime >= endTime {
		return Constraint{}, ErrInvalidArguments
	}
	return newConstraint(scheduleId, workerId, -1, -1, startTime, endTime, constraintType), nil
}

func NewTaskTimeConstraint(scheduleId, taskId, startTime, endTime int, constraintType ConstraintType) (Constraint, error) {
	if scheduleId < 0 || taskId < 0 || startTime >= endTime {
		return Constraint{}, ErrInvalidArguments
	}
	return newConstraint(scheduleId, -1, taskId, -1, startTime, endTime, constraintType), nil
}

func (c Constraint) ConflictsWith(other Constraint) bool {
	if c.Equals(other) {
		return true
	}

	if c.scheduleId != other.scheduleId {
		return false
	}

	if c.workerId != -1 && c.workerId == other.workerId && c.taskId != -1 && c.taskId == other.taskId {
		return true
	}

	if c.locationId != -1 && c.locationId == other.workerId && c.taskId != -1 && c.taskId == other.taskId {
		return true
	}

	if c.workerId != -1 && c.workerId == other.workerId && c.locationId != -1 && c.locationId == other.locationId {
		return true
	}

	if (c.locationId != -1 && c.locationId == other.locationId || c.taskId != -1 && c.taskId == other.taskId || c.workerId != -1 && c.workerId == other.workerId) &&
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
