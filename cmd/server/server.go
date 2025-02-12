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
	"github.com/pavozayac/constraints/config"
	"github.com/pavozayac/constraints/internal/application/protobuf"
	"github.com/pavozayac/constraints/internal/application/rpc"
	"github.com/pavozayac/constraints/internal/application/services"
	"github.com/pavozayac/constraints/internal/infrastructure/adapters"
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

	scheduleServer := rpc.NewScheduleServer(scheduleService)
	locationServer := rpc.NewLocationServer(locationService)
	taskServer := rpc.NewTaskServer(taskService)
	workerServer := rpc.NewWorkerServer(workerService)

	protobuf.RegisterScheduleServiceServer(s, scheduleServer)
	protobuf.RegisterLocationServiceServer(s, locationServer)
	protobuf.RegisterTaskServiceServer(s, taskServer)
	protobuf.RegisterWorkerServiceServer(s, workerServer)

	http.HandleFunc("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusOK)

		if _, err := io.WriteString(w, "OK"); err != nil {
			log.Fatalf("/health handler failed to write response: %v", err)
		}
	}))

	go func() {
		if err := http.ListenAndServe(":"+strconv.Itoa(conf.HeartbeatPort), nil); err != nil {
			log.Fatalf("failed to listen on port %d: %v", conf.HeartbeatPort, err)
		}
	}()

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
