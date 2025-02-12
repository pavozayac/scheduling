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
)

func TestLocationServer(t *testing.T) {
	ctx := context.Background()
	connStr, ctr, teardown := shared.SetupMigratedPostgresContainer(t, ctx, "../../infrastructure/db/migrations")
	defer teardown(t)

	err := ctr.Snapshot(context.Background())
	require.NoError(t, err)

	exampleTitle := "Location A"
	exampleDescription := "Description A"

	uuidOne := ishared.MockIdentityGenerator{}.Generate()
	uuidTwo := ishared.MockIdentityGenerator{}.Generate()
	uuidOneString := uuidOne.String()
	uuidTwoString := uuidTwo.String()

	t.Run("ShouldCreateNewLocation", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterLocationServiceServer(s, NewLocationServer(services.NewLocationService(adapters.NewPsqlLocationRepo(*pgConn))))
		}))
		locationClient := protobuf.NewLocationServiceClient(conn)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, "Test Schedule")
		require.NoError(t, err)

		request := &protobuf.CreateLocationRequest{
			ScheduleId:  &uuidOneString,
			Title:       &exampleTitle,
			Description: &exampleDescription,
		}

		response, err := locationClient.CreateLocation(context.Background(), request)
		assert.NoError(t, err)

		assert.NotNil(t, response.Id)
		assert.NotNil(t, uuidOneString, *response.ScheduleId)
		assert.Equal(t, exampleTitle, *response.Title)
		assert.Equal(t, exampleDescription, *response.Description)
	})

	t.Run("ShouldReadLocation", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterLocationServiceServer(s, NewLocationServer(services.NewLocationService(adapters.NewPsqlLocationRepo(*pgConn))))
		}))
		locationClient := protobuf.NewLocationServiceClient(conn)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, "Test Schedule")
		require.NoError(t, err)
		_, err = pgConn.Exec(context.Background(), "INSERT INTO locations (id, schedule_id, title, story) VALUES ($1, $2, $3, $4)", uuidTwo, uuidOne, exampleTitle, exampleDescription)
		require.NoError(t, err)

		request := &protobuf.LocationRequest{
			Id: &uuidTwoString,
		}

		response, err := locationClient.ReadLocation(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, uuidTwoString, *response.Id)
		assert.Equal(t, uuidOneString, *response.ScheduleId)
		assert.Equal(t, exampleTitle, *response.Title)
		assert.Equal(t, exampleDescription, *response.Description)
	})

	t.Run("ShouldDeleteLocation", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		conn := shared.SetupTestGrpcServer(t, shared.WithFunc(func(s *grpc.Server) {
			protobuf.RegisterLocationServiceServer(s, NewLocationServer(services.NewLocationService(adapters.NewPsqlLocationRepo(*pgConn))))
		}))
		locationClient := protobuf.NewLocationServiceClient(conn)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", uuidOne, "Test Schedule")
		require.NoError(t, err)
		_, err = pgConn.Exec(context.Background(), "INSERT INTO locations (id, schedule_id, title, story) VALUES ($1, $2, $3, $4)", uuidTwo, uuidOne, exampleTitle, exampleDescription)
		require.NoError(t, err)

		request := &protobuf.LocationRequest{
			Id: &uuidTwoString,
		}

		response, err := locationClient.DeleteLocation(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		var count int
		queryErr := pgConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM locations WHERE id = $1", uuidTwo).Scan(&count)
		assert.NoError(t, queryErr)
		assert.Equal(t, 0, count)
	})
}
