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
)

func TestTaskRepo(t *testing.T) {
	ctx := context.Background()
	ctr, teardown := tshared.SetupMigratedPostgresContainer(t, ctx)
	defer teardown(t)

	err := ctr.Snapshot(ctx)
	assert.NoError(t, err)

	connStr, err := ctr.ConnectionString(ctx)
	assert.NoError(t, err)

	t.Run("ShouldCreateNewTask", func(t *testing.T) {
		t.Cleanup(func() {
			err = ctr.Restore(ctx)
			assert.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		assert.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlTaskRepo{
			db: *pgConn,
		}

		taskId := shared.MockIdentityGenerator{}.Generate()
		scheduleId := shared.MockIdentityGenerator{}.Generate()

		task, err := model.NewTask(
			taskId,
			scheduleId,
			"Test Task",
			"Test description",
		)
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", scheduleId, "Test Schedule")
		assert.NoError(t, err)

		err = repo.SaveOrUpdateTask(context.Background(), *task)
		assert.NoError(t, err)

		var dbTaskId shared.Identity
		var dbScheduleId shared.Identity
		var dbTitle string
		var dbStory string

		queryErr := pgConn.QueryRow(context.Background(), "SELECT id, schedule_id, title, story FROM tasks WHERE id = $1", task.Id()).
			Scan(&dbTaskId, &dbScheduleId, &dbTitle, &dbStory)
		assert.NoError(t, err)
		assert.NoError(t, queryErr)

		assert.Equal(t, taskId, dbTaskId)
	})

	t.Run("ShouldRetrieveTask", func(t *testing.T) {
		t.Cleanup(func() {
			err = ctr.Restore(ctx)
			assert.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		assert.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlTaskRepo{
			db: *pgConn,
		}

		taskId := shared.MockIdentityGenerator{}.Generate()
		scheduleId := shared.MockIdentityGenerator{}.Generate()

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", scheduleId, "Test Schedule")
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO tasks (id, schedule_id, title, story) VALUES ($1, $2, $3, $4)", taskId, scheduleId, "Test Task", "Test description")
		assert.NoError(t, err)

		task, err := repo.GetTask(context.Background(), taskId)
		assert.NoError(t, err)

		assert.Equal(t, taskId, task.Id())
		assert.Equal(t, scheduleId, task.ScheduleId())
		assert.Equal(t, "Test Task", task.Name())
		assert.Equal(t, "Test description", task.Description())
	})

	t.Run("ShouldDeleteTask", func(t *testing.T) {
		t.Cleanup(func() {
			err = ctr.Restore(ctx)
			assert.NoError(t, err)
		})

		pgConn, err := pgx.Connect(context.Background(), connStr)
		assert.NoError(t, err)
		defer pgConn.Close(context.Background())

		repo := PsqlTaskRepo{
			db: *pgConn,
		}

		taskId := shared.MockIdentityGenerator{}.Generate()
		scheduleId := shared.MockIdentityGenerator{}.Generate()

		_, err = pgConn.Exec(context.Background(), "INSERT INTO schedules (id, title) VALUES ($1, $2)", scheduleId, "Test Schedule")
		assert.NoError(t, err)

		_, err = pgConn.Exec(context.Background(), "INSERT INTO tasks (id, schedule_id, title, story) VALUES ($1, $2, $3, $4)", taskId, scheduleId, "Test Task", "Test description")
		assert.NoError(t, err)

		err = repo.DeleteTask(context.Background(), taskId)
		assert.NoError(t, err)

		var count int
		queryErr := pgConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM tasks WHERE id = $1", taskId).Scan(&count)
		assert.NoError(t, queryErr)
		assert.Equal(t, 0, count)
	})
}
