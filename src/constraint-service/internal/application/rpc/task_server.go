package rpc

import (
	"context"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/protobuf"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
)

type TaskServer struct {
	taskService ports.TaskService

	protobuf.UnimplementedTaskServiceServer
}

func NewTaskServer(taskService ports.TaskService) *TaskServer {
	return &TaskServer{taskService: taskService}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *protobuf.CreateTaskRequest) (*protobuf.TaskResponse, error) {
	taskDTO := ports.TaskDTO{
		Id:          "",
		ScheduleId:  *req.ScheduleId,
		Title:       *req.Title,
		Description: *req.Description,
	}

	result, err := s.taskService.CreateOrModifyTask(ctx, taskDTO)
	if err != nil {
		return nil, err
	}

	return convertTaskDTOToProto(*result), nil
}

func (s *TaskServer) ReadTask(ctx context.Context, req *protobuf.TaskRequest) (*protobuf.TaskResponse, error) {
	task, err := s.taskService.GetTask(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return convertTaskDTOToProto(*task), nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *protobuf.TaskRequest) (*protobuf.TaskResponse, error) {
	err := s.taskService.RemoveTask(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return &protobuf.TaskResponse{}, nil
}
