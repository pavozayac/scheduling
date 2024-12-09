package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pavozayac/scheduling/src/constraint-service/config"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/protobuf"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/rpc"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/services"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/infrastructure/adapters"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf := config.LoadConfig()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Fprintf(os.Stderr, "%s", conf.DbConnectionString)

	db, err := pgx.Connect(context.Background(), conf.DbConnectionString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close(context.Background())

	if conf.UseStartupMigrate {
		gooseDb, err := goose.OpenDBWithDriver("pgx", conf.DbConnectionString)
		if err != nil {
			log.Fatalf("goose: failed to connect to database: %v", err)
		}
		defer gooseDb.Close()

		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatalf("goose: failed to set dialect: %v", err)
		}

		if err := goose.Up(gooseDb, conf.MigrationDirectory); err != nil {
			log.Fatalf("goose: failed to set dialect: %v", err)
		}
	}

	scheduleService := services.NewScheduleService(adapters.NewPsqlScheduleRepo(*db))
	locationService := services.NewLocationService(adapters.NewPsqlLocationRepo(*db))
	taskService := services.NewTaskService(adapters.NewPsqlTaskRepo(*db))
	workerService := services.NewWorkerService(adapters.NewPsqlWorkerRepo(*db))

	s := grpc.NewServer()

	grpcServer := rpc.NewGrpcServer(scheduleService, locationService, taskService, workerService)

	protobuf.RegisterScheduleServiceServer(s, &grpcServer)
	protobuf.RegisterLocationServiceServer(s, &grpcServer)
	protobuf.RegisterTaskServiceServer(s, &grpcServer)
	protobuf.RegisterWorkerServiceServer(s, &grpcServer)

	http.HandleFunc("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
	}))

	go http.ListenAndServe(":"+strconv.Itoa(conf.HeartbeatPort), nil)

	if conf.UseReflection {
		reflection.Register(s)
	}

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(conf.GrpcPort))

	if err != nil {
		log.Fatalf("failed to listen for TCP: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
