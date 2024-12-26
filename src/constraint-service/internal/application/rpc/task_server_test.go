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

func TestTaskServer(t *testing.T) {
	ctx := context.Background()
	connStr, ctr, teardown := shared.SetupMigratedPostgresContainer(t, ctx, "../../infrastructure/db/migrations")
	defer teardown(t)

	err := ctr.Snapshot(context.Background())
	require.NoError(t, err)

	exampleTitle := "Task A"
	exampleDescription := "Description A"

	uuidOne := ishared.MockIdentityGenerator{}.Generate()
	uuidTwo := ishared.MockIdentityGenerator{}.Generate()
	uuidOneString := uuidOne.String()
	uuidTwoString := uuidTwo.String()

	t.Run("ShouldCreateNewTask", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterTaskServiceServer(s, NewTaskServer(services.NewTaskService(adapters.NewPsqlTaskRepo(*pgConn))))
		}))
		taskClient := protobuf.NewTaskServiceClient(conn)

		pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, "Test Schedule")

		request := &protobuf.CreateTaskRequest{
			ScheduleId:  &uuidOneString,
			Title:       &exampleTitle,
			Description: &exampleDescription,
		}

		response, err := taskClient.CreateTask(context.Background(), request)
		assert.NoError(t, err)

		assert.NotNil(t, response.Id)
		assert.NotNil(t, uuidOneString, *response.ScheduleId)
		assert.Equal(t, exampleTitle, *response.Title)
		assert.Equal(t, exampleDescription, *response.Description)
	})

	t.Run("ShouldReadTask", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterTaskServiceServer(s, NewTaskServer(services.NewTaskService(adapters.NewPsqlTaskRepo(*pgConn))))
		}))
		taskClient := protobuf.NewTaskServiceClient(conn)

		pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, "Test Schedule")
		pgConn.Exec(context.Background(), "INSERT INTO tasks (id, schedule_id, title, story) VALUES ($1, $2, $3, $4)", uuidTwo, uuidOne, exampleTitle, exampleDescription)

		request := &protobuf.TaskRequest{
			Id: &uuidTwoString,
		}

		response, err := taskClient.ReadTask(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, uuidTwoString, *response.Id)
		assert.Equal(t, uuidOneString, *response.ScheduleId)
		assert.Equal(t, exampleTitle, *response.Title)
		assert.Equal(t, exampleDescription, *response.Description)
	})

	t.Run("ShouldDeleteTask", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterTaskServiceServer(s, NewTaskServer(services.NewTaskService(adapters.NewPsqlTaskRepo(*pgConn))))
		}))
		taskClient := protobuf.NewTaskServiceClient(conn)

		pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, "Test Schedule")
		pgConn.Exec(context.Background(), "INSERT INTO tasks (id, schedule_id, title, story) VALUES ($1, $2, $3, $4)", uuidTwo, uuidOne, exampleTitle, exampleDescription)

		request := &protobuf.TaskRequest{
			Id: &uuidTwoString,
		}

		response, err := taskClient.DeleteTask(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		var count int
		queryErr := pgConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM tasks WHERE id = $1", uuidTwo).Scan(&count)
		assert.NoError(t, queryErr)
		assert.Equal(t, 0, count)
	})
}
