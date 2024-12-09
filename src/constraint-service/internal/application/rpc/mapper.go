package rpc

import (
	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/protobuf"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
)

func convertConstraints(constraints []*protobuf.Constraint) []ports.ConstraintDTO {
	var result []ports.ConstraintDTO
	for _, c := range constraints {
		result = append(result, ports.ConstraintDTO{
			WorkerId:       *c.WorkerId,
			TaskId:         *c.TaskId,
			LocationId:     *c.LocationId,
			StartTime:      int(*c.StartTime),
			EndTime:        int(*c.EndTime),
			ConstraintType: c.Type.String(),
		})
	}
	return result
}

func convertScheduleDTOToProto(schedule ports.ScheduleDTO) *protobuf.ScheduleResponse {
	return &protobuf.ScheduleResponse{
		Id:          &schedule.Id,
		Title:       &schedule.Title,
		Constraints: convertConstraintsToProto(schedule.Constraints),
	}
}

func convertConstraintsToProto(constraints []ports.ConstraintDTO) []*protobuf.Constraint {
	var result []*protobuf.Constraint
	for _, c := range constraints {
		intStartTime := int32(c.StartTime)
		intEndTime := int32(c.EndTime)

		var protoConstraintType protobuf.ConstraintType
		switch c.ConstraintType {
		case "must":
			protoConstraintType = protobuf.ConstraintType_MUST
		case "cannot":
			protoConstraintType = protobuf.ConstraintType_CANNOT
		default:
			protoConstraintType = protobuf.ConstraintType_MUST
		}

		result = append(result, &protobuf.Constraint{
			WorkerId:   &c.WorkerId,
			TaskId:     &c.TaskId,
			LocationId: &c.LocationId,
			StartTime:  &intStartTime,
			EndTime:    &intEndTime,
			Type:       &protoConstraintType,
		})
	}
	return result
}

func convertLocationDTOToProto(location ports.LocationDTO) *protobuf.LocationResponse {
	return &protobuf.LocationResponse{
		Id:          &location.Id,
		ScheduleId:  &location.ScheduleId,
		Title:       &location.Name,
		Description: &location.Description,
	}
}

func convertTaskDTOToProto(task ports.TaskDTO) *protobuf.TaskResponse {
	return &protobuf.TaskResponse{
		Id:          &task.Id,
		ScheduleId:  &task.ScheduleId,
		Title:       &task.Title,
		Description: &task.Description,
	}
}

func convertWorkerDTOToProto(worker ports.WorkerDTO) *protobuf.WorkerResponse {
	return &protobuf.WorkerResponse{
		Id:         &worker.Id,
		ScheduleId: &worker.ScheduleId,
		FirstName:  &worker.FirstName,
		LastName:   &worker.LastName,
	}
}
