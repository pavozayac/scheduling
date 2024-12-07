package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/protobuf"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/rpc"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/services"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/infrastructure/adapters"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	s := grpc.NewServer()

	db, err := pgx.Connect(context.Background(), "postgresql://postgres:password@localhost:5432/scheduling")

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	scheduleService := services.NewScheduleService(adapters.NewPsqlScheduleRepo(*db))
	locationService := services.NewLocationService(adapters.NewPsqlLocationRepo(*db))
	taskService := services.NewTaskService(adapters.NewPsqlTaskRepo(*db))
	workerService := services.NewWorkerService(adapters.NewPsqlWorkerRepo(*db))

	grpcServer := rpc.NewGrpcServer(scheduleService, locationService, taskService, workerService)

	protobuf.RegisterScheduleServiceServer(s, &grpcServer)
	protobuf.RegisterLocationServiceServer(s, &grpcServer)
	protobuf.RegisterTaskServiceServer(s, &grpcServer)
	protobuf.RegisterWorkerServiceServer(s, &grpcServer)

	reflection.Register(s)

	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("failed to listen for TCP: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
