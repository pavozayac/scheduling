package adapters

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"

	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/infrastructure/db/sqlc"
	ishared "github.com/pavozayac/scheduling/src/constraint-service/internal/infrastructure/shared"
)

type PsqlTaskRepo struct {
	db pgx.Conn
}

func (r PsqlTaskRepo) Db() pgx.Conn {
	return r.db
}

func NewPsqlTaskRepo(db pgx.Conn) *PsqlTaskRepo {
	return &PsqlTaskRepo{db: db}
}

func (r *PsqlTaskRepo) SaveOrUpdateTask(ctx context.Context, task model.Task) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		taskID := uuid.UUID(task.Id())
		scheduleID := uuid.UUID(task.ScheduleId())
		return q.UpsertTask(ctx, sqlc.UpsertTaskParams{
			ID:         &taskID,
			Title:      task.Name(),
			Story:      task.Description(),
			ScheduleID: &scheduleID,
		})
	})
}

func (r *PsqlTaskRepo) GetTask(ctx context.Context, id shared.Identity) (task *model.Task, err error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = errors.Join(err, tx.Rollback(ctx))
	}()

	queries := sqlc.New(tx)

	taskID := uuid.UUID(id)
	dbTask, err := queries.GetTask(ctx, &taskID)
	if err != nil {
		return nil, err
	}

	return model.NewTask(
		shared.Identity(*dbTask.ID),
		shared.Identity(*dbTask.ScheduleID),
		dbTask.Title,
		dbTask.Story,
	)
}

func (r *PsqlTaskRepo) DeleteTask(ctx context.Context, id shared.Identity) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		taskID := uuid.UUID(id)
		return q.DeleteTask(ctx, &taskID)
	})
}

var _ ports.TaskRepository = &PsqlTaskRepo{}
