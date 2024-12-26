package rpc

import (
	"context"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/protobuf"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
)

type LocationServer struct {
	locationService ports.LocationService

	protobuf.UnimplementedLocationServiceServer
}

func NewLocationServer(locationService ports.LocationService) *LocationServer {
	return &LocationServer{locationService: locationService}
}

func (s *LocationServer) CreateLocation(ctx context.Context, req *protobuf.CreateLocationRequest) (*protobuf.LocationResponse, error) {
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

func (s *LocationServer) ReadLocation(ctx context.Context, req *protobuf.LocationRequest) (*protobuf.LocationResponse, error) {
	location, err := s.locationService.GetLocation(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return convertLocationDTOToProto(*location), nil
}

func (s *LocationServer) DeleteLocation(ctx context.Context, req *protobuf.LocationRequest) (*protobuf.LocationResponse, error) {
	err := s.locationService.RemoveLocation(ctx, *req.Id)
	if err != nil {
		return nil, err
	}
	return &protobuf.LocationResponse{}, nil
}
