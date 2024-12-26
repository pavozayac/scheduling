package rpc

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/protobuf"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/services"
	ishared "github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/infrastructure/adapters"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/tests/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestWorkerServer(t *testing.T) {
	ctx := context.Background()
	connStr, ctr, teardown := shared.SetupMigratedPostgresContainer(t, ctx, "../../infrastructure/db/migrations")
	defer teardown(t)

	err := ctr.Snapshot(context.Background())
	require.NoError(t, err)

	exampleFirstName := "John"
	exampleLastName := "Doe"

	uuidOne := ishared.MockIdentityGenerator{}.Generate()
	uuidTwo := ishared.MockIdentityGenerator{}.Generate()
	uuidOneString := uuidOne.String()
	uuidTwoString := uuidTwo.String()

	t.Run("ShouldCreateNewWorker", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterWorkerServiceServer(s, NewWorkerServer(services.NewWorkerService(adapters.NewPsqlWorkerRepo(*pgConn))))
		}))
		workerClient := protobuf.NewWorkerServiceClient(conn)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, "Test Schedule")
		require.NoError(t, err)

		request := &protobuf.CreateWorkerRequest{
			ScheduleId: &uuidOneString,
			FirstName:  &exampleFirstName,
			LastName:   &exampleLastName,
		}

		response, err := workerClient.CreateWorker(context.Background(), request)
		assert.NoError(t, err)

		assert.NotNil(t, response.Id)
		assert.NotNil(t, uuidOneString, *response.ScheduleId)
		assert.Equal(t, exampleFirstName, *response.FirstName)
		assert.Equal(t, exampleLastName, *response.LastName)
	})

	t.Run("ShouldReadWorker", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterWorkerServiceServer(s, NewWorkerServer(services.NewWorkerService(adapters.NewPsqlWorkerRepo(*pgConn))))
		}))
		workerClient := protobuf.NewWorkerServiceClient(conn)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, "Test Schedule")
		require.NoError(t, err)
		_, err = pgConn.Exec(context.Background(), "INSERT INTO workers (id, schedule_id, first_name, last_name) VALUES ($1, $2, $3, $4)", uuidTwo, uuidOne, exampleFirstName, exampleLastName)
		require.NoError(t, err)

		request := &protobuf.WorkerRequest{
			Id: &uuidTwoString,
		}

		response, err := workerClient.ReadWorker(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, uuidTwoString, *response.Id)
		assert.Equal(t, uuidOneString, *response.ScheduleId)
		assert.Equal(t, exampleFirstName, *response.FirstName)
		assert.Equal(t, exampleLastName, *response.LastName)
	})

	t.Run("ShouldDeleteWorker", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterWorkerServiceServer(s, NewWorkerServer(services.NewWorkerService(adapters.NewPsqlWorkerRepo(*pgConn))))
		}))
		workerClient := protobuf.NewWorkerServiceClient(conn)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, "Test Schedule")
		require.NoError(t, err)
		_, err = pgConn.Exec(context.Background(), "INSERT INTO workers (id, schedule_id, first_name, last_name) VALUES ($1, $2, $3, $4)", uuidTwo, uuidOne, exampleFirstName, exampleLastName)
		require.NoError(t, err)

		request := &protobuf.WorkerRequest{
			Id: &uuidTwoString,
		}

		response, err := workerClient.DeleteWorker(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		var count int
		queryErr := pgConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM workers WHERE id = $1", uuidTwo).Scan(&count)
		assert.NoError(t, queryErr)
		assert.Equal(t, 0, count)
	})
}
