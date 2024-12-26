package rpc

import (
	"context"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/protobuf"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
)

type WorkerServer struct {
	workerService ports.WorkerService

	protobuf.UnimplementedWorkerServiceServer
}

func NewWorkerServer(workerService ports.WorkerService) *WorkerServer {
	return &WorkerServer{workerService: workerService}
}

func (s *WorkerServer) CreateWorker(ctx context.Context, req *protobuf.CreateWorkerRequest) (*protobuf.WorkerResponse, error) {
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

func (s *WorkerServer) ReadWorker(ctx context.Context, req *protobuf.WorkerRequest) (*protobuf.WorkerResponse, error) {
	worker, err := s.workerService.GetWorker(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return convertWorkerDTOToProto(*worker), nil
}

func (s *WorkerServer) DeleteWorker(ctx context.Context, req *protobuf.WorkerRequest) (*protobuf.WorkerResponse, error) {
	err := s.workerService.RemoveWorker(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return &protobuf.WorkerResponse{}, nil
}
