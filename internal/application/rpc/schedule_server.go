package rpc

import (
	"context"
	"fmt"

	"github.com/pavozayac/constraints/internal/application/protobuf"
	"github.com/pavozayac/constraints/internal/domain/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ScheduleServer struct {
	scheduleService ports.ScheduleService

	protobuf.UnimplementedScheduleServiceServer
}

func NewScheduleServer(scheduleService ports.ScheduleService) *ScheduleServer {
	return &ScheduleServer{scheduleService: scheduleService}
}

func (s *ScheduleServer) CreateSchedule(ctx context.Context, req *protobuf.CreateScheduleRequest) (*protobuf.ScheduleResponse, error) {
	if req.Title == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}

	scheduleDTO := ports.ScheduleDTO{
		Id:          "",
		Title:       *req.Title,
		Constraints: convertConstraints(req.Constraints),
	}

	result, err := s.scheduleService.CreateOrModifySchedule(ctx, scheduleDTO)

	if err != nil {
		return nil, fmt.Errorf("scheduleService.CreateOrModifySchedule() error: %w", err)
	}

	return convertScheduleDTOToProto(*result), nil
}

func (s *ScheduleServer) ReadSchedule(ctx context.Context, req *protobuf.ScheduleRequest) (*protobuf.ScheduleResponse, error) {
	schedule, err := s.scheduleService.GetSchedule(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return convertScheduleDTOToProto(*schedule), nil
}

func (s *ScheduleServer) UpdateSchedule(ctx context.Context, req *protobuf.UpdateScheduleRequest) (*protobuf.ScheduleResponse, error) {
	if req.Title == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing required fields")
	}

	scheduleDTO := ports.ScheduleDTO{
		Id:          *req.Id,
		Title:       *req.Title,
		Constraints: convertConstraints(req.Constraints),
	}

	result, err := s.scheduleService.CreateOrModifySchedule(ctx, scheduleDTO)

	if err != nil {
		return nil, fmt.Errorf("scheduleService.CreateOrModifySchedule() error: %w", err)
	}

	return convertScheduleDTOToProto(*result), nil
}

func (s *ScheduleServer) DeleteSchedule(ctx context.Context, req *protobuf.ScheduleRequest) (*protobuf.ScheduleResponse, error) {
	err := s.scheduleService.RemoveSchedule(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return &protobuf.ScheduleResponse{}, nil
}
