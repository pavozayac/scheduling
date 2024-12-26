package adapters

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	tshared "github.com/pavozayac/scheduling/src/constraint-service/internal/tests/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScheduleRepo(t *testing.T) {
	ctx := context.Background()
	connStr, ctr, teardown := tshared.SetupMigratedPostgresContainer(t, ctx, "../db/migrations")
	defer teardown(t)

	err := ctr.Snapshot(context.Background())
	require.NoError(t, err)

	t.Run("ShouldCreateNewSchedule", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlScheduleRepo{
			db: *pgConn,
		}

		scheduleId := shared.MockIdentityGenerator{}.Generate()

		schedule, err := model.NewSchedule(
			scheduleId,
			"Test Schedule",
			[]model.Constraint{},
		)
		assert.NoError(t, err)

		err = repo.SaveOrUpdateSchedule(context.Background(), *schedule)
		assert.NoError(t, err)

		var dbScheduleId shared.Identity
		var dbTitle string
		queryErr := pgConn.QueryRow(context.Background(),
			"SELECT id, title FROM schedules WHERE id = $1",
			scheduleId).Scan(&dbScheduleId, &dbTitle)
		assert.NoError(t, queryErr)
		assert.Equal(t, scheduleId, dbScheduleId)
		assert.Equal(t, "Test Schedule", dbTitle)
	})

	t.Run("ShouldRetrieveSchedule", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlScheduleRepo{
			db: *pgConn,
		}

		scheduleId := shared.MockIdentityGenerator{}.Generate()
		workerId := shared.MockIdentityGenerator{}.Generate()
		taskId := shared.NilIdentity
		locationId := shared.NilIdentity

		_, err = pgConn.Exec(context.Background(),
			"INSERT INTO schedules (id, title) VALUES ($1, $2)",
			scheduleId, "Test Schedule")
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(),
			"INSERT INTO workers (id, first_name, last_name, schedule_id) VALUES ($1, $2, $3, $4)",
			workerId, "Test Name", "Test Last Name", scheduleId)
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(),
			`INSERT INTO constraints 
            (schedule_id, worker_id, start_slot, end_slot, kind) 
            VALUES ($1, $2, $3, $4, $5)`,
			scheduleId, workerId, 1, 2, "must")
		assert.NoError(t, err)

		schedule, err := repo.GetSchedule(context.Background(), scheduleId)
		assert.NoError(t, err)
		assert.Equal(t, scheduleId, schedule.Id())
		assert.Equal(t, "Test Schedule", schedule.Title())
		assert.Equal(t, 1, len(schedule.Constraints()))

		constraint := schedule.Constraints()[0]
		assert.Equal(t, scheduleId, constraint.ScheduleId())
		assert.Equal(t, workerId, constraint.WorkerId())
		assert.Equal(t, taskId, constraint.TaskId())
		assert.Equal(t, locationId, constraint.LocationId())
		assert.Equal(t, 1, constraint.StartTime())
		assert.Equal(t, 2, constraint.EndTime())
		assert.Equal(t, model.Must, constraint.Type())
	})

	t.Run("ShouldUpdateSchedule", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			assert.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlScheduleRepo{
			db: *pgConn,
		}

		scheduleId := shared.MockIdentityGenerator{}.Generate()
		workerId := shared.MockIdentityGenerator{}.Generate()
		taskId := shared.NilIdentity
		locationId := shared.NilIdentity

		_, err = pgConn.Exec(context.Background(),
			"INSERT INTO schedules (id, title) VALUES ($1, $2)",
			scheduleId, "Test Schedule")
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(),
			"INSERT INTO workers (id, first_name, last_name, schedule_id) VALUES ($1, $2, $3, $4)",
			workerId, "Test Name", "Test Last Name", scheduleId)
		assert.NoError(t, err)

		constraint, err := model.NewConstraint(
			scheduleId,
			workerId,
			taskId,
			locationId,
			3,
			4,
			model.Must,
		)
		assert.NoError(t, err)

		schedule, err := model.NewSchedule(
			scheduleId,
			"Updated Schedule",
			[]model.Constraint{constraint},
		)
		assert.NoError(t, err)

		err = repo.SaveOrUpdateSchedule(context.Background(), *schedule)
		assert.NoError(t, err)

		var dbTitle string
		queryErr := pgConn.QueryRow(context.Background(),
			"SELECT title FROM schedules WHERE id = $1",
			scheduleId).Scan(&dbTitle)
		assert.NoError(t, queryErr)
		assert.Equal(t, "Updated Schedule", dbTitle)

		var startSlot, endSlot int
		queryErr = pgConn.QueryRow(context.Background(),
			"SELECT start_slot, end_slot FROM constraints WHERE schedule_id = $1",
			scheduleId).Scan(&startSlot, &endSlot)
		assert.NoError(t, queryErr)
		assert.Equal(t, 3, startSlot)
		assert.Equal(t, 4, endSlot)
	})

	t.Run("ShouldDeleteSchedule", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlScheduleRepo{
			db: *pgConn,
		}

		scheduleId := shared.MockIdentityGenerator{}.Generate()
		workerId := shared.MockIdentityGenerator{}.Generate()

		_, err = pgConn.Exec(context.Background(),
			"INSERT INTO schedules (id, title) VALUES ($1, $2)",
			scheduleId, "Test Schedule")
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(),
			"INSERT INTO workers (id, first_name, last_name, schedule_id) VALUES ($1, $2, $3, $4)",
			workerId, "Test Name", "Test Last Name", scheduleId)
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(),
			`INSERT INTO constraints 
            (schedule_id, worker_id, start_slot, end_slot, kind) 
            VALUES ($1, $2, $3, $4, $5)`,
			scheduleId, workerId, 1, 2, "must")
		assert.NoError(t, err)

		err = repo.DeleteSchedule(context.Background(), scheduleId)
		assert.NoError(t, err)

		var count int
		queryErr := pgConn.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM schedules WHERE id = $1",
			scheduleId).Scan(&count)
		assert.NoError(t, queryErr)
		assert.Equal(t, 0, count)

		queryErr = pgConn.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM constraints WHERE schedule_id = $1",
			scheduleId).Scan(&count)
		assert.NoError(t, queryErr)
		assert.Equal(t, 0, count)
	})
}
