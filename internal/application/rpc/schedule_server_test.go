package rpc

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pavozayac/constraints/internal/application/protobuf"
	"github.com/pavozayac/constraints/internal/application/services"
	ishared "github.com/pavozayac/constraints/internal/domain/shared"
	"github.com/pavozayac/constraints/internal/infrastructure/adapters"
	"github.com/pavozayac/constraints/internal/tests/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func TestScheduleServer(t *testing.T) {
	ctx := context.Background()
	connStr, ctr, teardown := shared.SetupMigratedPostgresContainer(t, ctx, "../../infrastructure/db/migrations")
	defer teardown(t)

	err := ctr.Snapshot(context.Background())
	require.NoError(t, err)

	uuidOne := ishared.MockIdentityGenerator{}.Generate()
	uuidOneString := uuidOne.String()
	uuidTwo := ishared.MockIdentityGenerator{}.Generate()
	uuidTwoString := uuidTwo.String()

	exampleTitle := "Schedule A"

	exampleConstraint := &protobuf.Constraint{
		WorkerId:  proto.String(uuidTwoString),
		StartTime: proto.Int32(9),
		EndTime:   proto.Int32(17),
		Type:      protobuf.ConstraintType_must.Enum(),
	}

	t.Run("ShouldCreateNewSchedule", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterScheduleServiceServer(s, NewScheduleServer(services.NewScheduleService(adapters.NewPsqlScheduleRepo(*pgConn))))
		}))
		scheduleClient := protobuf.NewScheduleServiceClient(conn)

		request := &protobuf.CreateScheduleRequest{
			Title: &exampleTitle,
		}

		response, err := scheduleClient.CreateSchedule(context.Background(), request)
		assert.NoError(t, err)

		assert.NotNil(t, response.Id)
		assert.Equal(t, exampleTitle, *response.Title)
	})

	t.Run("ShouldReadSchedule", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterScheduleServiceServer(s, NewScheduleServer(services.NewScheduleService(adapters.NewPsqlScheduleRepo(*pgConn))))
		}))
		scheduleClient := protobuf.NewScheduleServiceClient(conn)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, exampleTitle)
		require.NoError(t, err)

		request := &protobuf.ScheduleRequest{
			Id: &uuidOneString,
		}

		response, err := scheduleClient.ReadSchedule(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, uuidOneString, *response.Id)
		assert.Equal(t, exampleTitle, *response.Title)
	})

	t.Run("ShouldDeleteSchedule", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterScheduleServiceServer(s, NewScheduleServer(services.NewScheduleService(adapters.NewPsqlScheduleRepo(*pgConn))))
		}))
		scheduleClient := protobuf.NewScheduleServiceClient(conn)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, exampleTitle)
		require.NoError(t, err)

		request := &protobuf.ScheduleRequest{
			Id: &uuidOneString,
		}

		response, err := scheduleClient.DeleteSchedule(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		var count int
		queryErr := pgConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM schedules WHERE id = $1", uuidOne).Scan(&count)
		assert.NoError(t, queryErr)
		assert.Equal(t, 0, count)
	})

	t.Run("ShouldAddConstraintToSchedule", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterScheduleServiceServer(s, NewScheduleServer(services.NewScheduleService(adapters.NewPsqlScheduleRepo(*pgConn))))
		}))
		scheduleClient := protobuf.NewScheduleServiceClient(conn)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, exampleTitle)
		require.NoError(t, err)
		_, err = pgConn.Exec(context.Background(), "INSERT INTO workers (id, schedule_id, first_name, last_name) VALUES ($1, $2, $3, $4)", uuidTwo, uuidOne, "John", "Doe")
		require.NoError(t, err)

		request := &protobuf.UpdateScheduleRequest{
			Id:          &uuidOneString,
			Title:       &exampleTitle,
			Constraints: []*protobuf.Constraint{exampleConstraint},
		}

		response, err := scheduleClient.UpdateSchedule(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		assert.Equal(t, uuidOneString, *response.Id)
		assert.Equal(t, exampleTitle, *response.Title)
		assert.Equal(t, 1, len(response.Constraints))
		assert.EqualValues(t, exampleConstraint.WorkerId, response.Constraints[0].WorkerId)
		assert.EqualValues(t, exampleConstraint.TaskId, response.Constraints[0].TaskId)
		assert.EqualValues(t, exampleConstraint.LocationId, response.Constraints[0].LocationId)
		assert.EqualValues(t, exampleConstraint.StartTime, response.Constraints[0].StartTime)
		assert.EqualValues(t, exampleConstraint.EndTime, response.Constraints[0].EndTime)
		assert.EqualValues(t, exampleConstraint.Type, response.Constraints[0].Type)

	})

	t.Run("ShouldRemoveConstraintFromSchedule", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterScheduleServiceServer(s, NewScheduleServer(services.NewScheduleService(adapters.NewPsqlScheduleRepo(*pgConn))))
		}))
		scheduleClient := protobuf.NewScheduleServiceClient(conn)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, exampleTitle)
		require.NoError(t, err)
		_, err = pgConn.Exec(context.Background(), "INSERT INTO workers (id, schedule_id, first_name, last_name) VALUES ($1, $2, $3, $4)", uuidTwo, uuidOne, "John", "Doe")
		require.NoError(t, err)
		_, err = pgConn.Exec(context.Background(), "INSERT INTO constraints (schedule_id, worker_id, task_id, location_id, start_slot, end_slot, kind) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			uuidOne, exampleConstraint.WorkerId, exampleConstraint.TaskId, exampleConstraint.LocationId, exampleConstraint.StartTime, exampleConstraint.EndTime, exampleConstraint.Type.String())
		require.NoError(t, err)

		request := &protobuf.UpdateScheduleRequest{
			Id:          &uuidOneString,
			Title:       &exampleTitle,
			Constraints: []*protobuf.Constraint{},
		}

		response, err := scheduleClient.UpdateSchedule(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		assert.Equal(t, uuidOneString, *response.Id)
		assert.Equal(t, exampleTitle, *response.Title)
		assert.Equal(t, 0, len(response.Constraints))
	})
}
