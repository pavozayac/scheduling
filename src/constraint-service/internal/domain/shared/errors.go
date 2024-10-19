package shared

import "errors"

var ErrNegativeId = errors.New("id arguments must be greater than or equal to 0")
var ErrInvalidArguments = errors.New("scheduleId, workerId, and taskId must not be nullish")
