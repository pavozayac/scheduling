package adapters

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pavozayac/constraints/internal/domain/model"
	"github.com/pavozayac/constraints/internal/domain/shared"
	tshared "github.com/pavozayac/constraints/internal/tests/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocationRepo(t *testing.T) {
	ctx := context.Background()
	connStr, ctr, teardown := tshared.SetupMigratedPostgresContainer(t, ctx, "../db/migrations")
	defer teardown(t)

	err := ctr.Snapshot(context.Background())
	require.NoError(t, err)

	t.Run("ShouldCreateNewLocation", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlLocationRepo{
			db: *pgConn,
		}

		locationId := shared.MockIdentityGenerator{}.Generate()
		scheduleId := shared.MockIdentityGenerator{}.Generate()

		location, err := model.NewLocation(
			locationId,
			scheduleId,
			"Test Location",
			"Test description",
		)
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", scheduleId, "Test Schedule")
		assert.NoError(t, err)

		err = repo.SaveOrUpdateLocation(context.Background(), *location)
		assert.NoError(t, err)

		var dbLocationId shared.Identity
		var dbScheduleId shared.Identity
		var dbTitle string
		var dbStory string

		queryErr := pgConn.QueryRow(context.Background(), "SELECT id, schedule_id, title, story FROM locations WHERE id = $1", location.Id()).
			Scan(&dbLocationId, &dbScheduleId, &dbTitle, &dbStory)
		assert.NoError(t, err)
		assert.NoError(t, queryErr)

		assert.Equal(t, locationId, dbLocationId)
	})

	t.Run("ShouldRetrieveLocation", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			assert.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlLocationRepo{
			db: *pgConn,
		}

		locationId := shared.MockIdentityGenerator{}.Generate()
		scheduleId := shared.MockIdentityGenerator{}.Generate()

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", scheduleId, "Test Schedule")
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO locations (id, schedule_id, title, story) VALUES ($1, $2, $3, $4)", locationId, scheduleId, "Test Location", "Test description")
		assert.NoError(t, err)

		location, err := repo.GetLocation(context.Background(), locationId)
		assert.NoError(t, err)

		assert.Equal(t, locationId, location.Id())
		assert.Equal(t, scheduleId, location.ScheduleId())
		assert.Equal(t, "Test Location", location.Name())
		assert.Equal(t, "Test description", location.Description())
	})

	t.Run("ShouldDeleteLocation", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlLocationRepo{
			db: *pgConn,
		}

		locationId := shared.MockIdentityGenerator{}.Generate()
		scheduleId := shared.MockIdentityGenerator{}.Generate()

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", scheduleId, "Test Schedule")
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO locations (id, schedule_id, title, story) VALUES ($1, $2, $3, $4)", locationId, scheduleId, "Test Location", "Test description")
		assert.NoError(t, err)

		err = repo.DeleteLocation(context.Background(), locationId)
		assert.NoError(t, err)

		var count int
		queryErr := pgConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM locations WHERE id = $1", locationId).Scan(&count)
		assert.NoError(t, queryErr)
		assert.Equal(t, 0, count)
	})
}
