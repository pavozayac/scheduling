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

func TestWorkerRepo(t *testing.T) {
	ctx := context.Background()
	connStr, ctr, teardown := tshared.SetupMigratedPostgresContainer(t, ctx, "../db/migrations")
	defer teardown(t)

	err := ctr.Snapshot(context.Background())
	require.NoError(t, err)

	t.Run("ShouldCreateNewWorker", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlWorkerRepo{
			db: *pgConn,
		}

		workerId := shared.MockIdentityGenerator{}.Generate()
		scheduleId := shared.MockIdentityGenerator{}.Generate()

		worker, err := model.NewWorker(
			workerId,
			scheduleId,
			"John",
			"Doe",
		)
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", scheduleId, "Test Schedule")
		assert.NoError(t, err)

		err = repo.SaveOrUpdateWorker(context.Background(), *worker)
		assert.NoError(t, err)

		var dbWorkerId shared.Identity
		var dbScheduleId shared.Identity
		var dbFirstName string
		var dbLastName string

		queryErr := pgConn.QueryRow(context.Background(), "SELECT id, schedule_id, first_name, last_name FROM workers WHERE id = $1", worker.Id()).
			Scan(&dbWorkerId, &dbScheduleId, &dbFirstName, &dbLastName)
		assert.NoError(t, queryErr)

		assert.Equal(t, workerId, dbWorkerId)
		assert.Equal(t, scheduleId, dbScheduleId)
		assert.Equal(t, "John", dbFirstName)
		assert.Equal(t, "Doe", dbLastName)
	})

	t.Run("ShouldRetrieveWorker", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlWorkerRepo{
			db: *pgConn,
		}

		workerId := shared.MockIdentityGenerator{}.Generate()
		scheduleId := shared.MockIdentityGenerator{}.Generate()

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", scheduleId, "Test Schedule")
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO workers (id, schedule_id, first_name, last_name) VALUES ($1, $2, $3, $4)",
			workerId, scheduleId, "John", "Doe")
		assert.NoError(t, err)

		worker, err := repo.GetWorker(context.Background(), workerId)
		assert.NoError(t, err)

		assert.Equal(t, workerId, worker.Id())
		assert.Equal(t, scheduleId, worker.ScheduleId())
		assert.Equal(t, "John", worker.FirstName())
		assert.Equal(t, "Doe", worker.LastName())
	})

	t.Run("ShouldDeleteWorker", func(t *testing.T) {
		t.Cleanup(func() {
			err := ctr.Restore(ctx)
			require.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		require.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlWorkerRepo{
			db: *pgConn,
		}

		workerId := shared.MockIdentityGenerator{}.Generate()
		scheduleId := shared.MockIdentityGenerator{}.Generate()

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", scheduleId, "Test Schedule")
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO workers (id, schedule_id, first_name, last_name) VALUES ($1, $2, $3, $4)",
			workerId, scheduleId, "John", "Doe")
		assert.NoError(t, err)

		err = repo.DeleteWorker(context.Background(), workerId)
		assert.NoError(t, err)

		var count int
		queryErr := pgConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM workers WHERE id = $1", workerId).Scan(&count)
		assert.NoError(t, queryErr)
		assert.Equal(t, 0, count)
	})
}
