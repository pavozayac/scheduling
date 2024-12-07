package rpc

import (
	"context"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/protobuf"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	scheduleService ports.ScheduleService
	locationService ports.LocationService
	taskService     ports.TaskService
	workerService   ports.WorkerService

	protobuf.UnimplementedScheduleServiceServer
	protobuf.UnimplementedLocationServiceServer
	protobuf.UnimplementedTaskServiceServer
	protobuf.UnimplementedWorkerServiceServer
}

var _ protobuf.ScheduleServiceServer = &GrpcServer{}
var _ protobuf.LocationServiceServer = &GrpcServer{}
var _ protobuf.TaskServiceServer = &GrpcServer{}
var _ protobuf.WorkerServiceServer = &GrpcServer{}

func NewGrpcServer(
	scheduleService ports.ScheduleService,
	locationService ports.LocationService,
	taskService ports.TaskService,
	workerService ports.WorkerService,
) GrpcServer {
	return GrpcServer{
		scheduleService: scheduleService,
		locationService: locationService,
		taskService:     taskService,
		workerService:   workerService,
	}
}

func (s *GrpcServer) CreateSchedule(ctx context.Context, req *protobuf.CreateScheduleRequest) (*protobuf.ScheduleResponse, error) {
	if req.Title == nil || req.Constraints == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}

	scheduleDTO := ports.ScheduleDTO{
		Id:          "",
		Title:       *req.Title,
		Constraints: convertConstraints(req.Constraints),
	}

	result, err := s.scheduleService.CreateOrModifySchedule(ctx, scheduleDTO)

	if err != nil {
		return nil, err
	}

	return convertScheduleDTOToProto(*result), nil
}

func (s *GrpcServer) ReadSchedule(ctx context.Context, req *protobuf.ScheduleRequest) (*protobuf.ScheduleResponse, error) {
	schedule, err := s.scheduleService.GetSchedule(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return convertScheduleDTOToProto(*schedule), nil
}

func (s *GrpcServer) DeleteSchedule(ctx context.Context, req *protobuf.ScheduleRequest) (*protobuf.ScheduleResponse, error) {
	err := s.scheduleService.RemoveSchedule(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return &protobuf.ScheduleResponse{}, nil
}

func (s *GrpcServer) CreateLocation(ctx context.Context, req *protobuf.CreateLocationRequest) (*protobuf.LocationResponse, error) {
	locationDTO := ports.LocationDTO{
		Id:          "",
		ScheduleId:  *req.ScheduleId,
		Name:        *req.Title,
		Description: *req.Description,
	}

	result, err := s.locationService.CreateOrModifyLocation(ctx, locationDTO)

	if err != nil {
		return nil, err
	}

	return convertLocationDTOToProto(*result), nil
}

func (s *GrpcServer) ReadLocation(ctx context.Context, req *protobuf.LocationRequest) (*protobuf.LocationResponse, error) {
	location, err := s.locationService.GetLocation(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return convertLocationDTOToProto(*location), nil
}

func (s *GrpcServer) DeleteLocation(ctx context.Context, req *protobuf.LocationRequest) (*protobuf.LocationResponse, error) {
	err := s.locationService.RemoveLocation(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return &protobuf.LocationResponse{}, nil
}

func (s *GrpcServer) CreateTask(ctx context.Context, req *protobuf.CreateTaskRequest) (*protobuf.TaskResponse, error) {
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

func (s *GrpcServer) ReadTask(ctx context.Context, req *protobuf.TaskRequest) (*protobuf.TaskResponse, error) {
	task, err := s.taskService.GetTask(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return convertTaskDTOToProto(*task), nil
}

func (s *GrpcServer) DeleteTask(ctx context.Context, req *protobuf.TaskRequest) (*protobuf.TaskResponse, error) {
	err := s.taskService.RemoveTask(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return &protobuf.TaskResponse{}, nil
}

func (s *GrpcServer) CreateWorker(ctx context.Context, req *protobuf.CreateWorkerRequest) (*protobuf.WorkerResponse, error) {
	workerDTO := ports.WorkerDTO{
		Id:         "",
		ScheduleId: *req.ScheduleId,
		FirstName:  *req.FirstName,
		LastName:   *req.LastName,
	}

	result, err := s.workerService.CreateOrModifyWorker(ctx, workerDTO)
	if err != nil {
		return nil, err
	}

	return convertWorkerDTOToProto(*result), nil
}

func (s *GrpcServer) ReadWorker(ctx context.Context, req *protobuf.WorkerRequest) (*protobuf.WorkerResponse, error) {
	worker, err := s.workerService.GetWorker(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return convertWorkerDTOToProto(*worker), nil
}

func (s *GrpcServer) DeleteWorker(ctx context.Context, req *protobuf.WorkerRequest) (*protobuf.WorkerResponse, error) {
	err := s.workerService.RemoveWorker(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return &protobuf.WorkerResponse{}, nil
}
