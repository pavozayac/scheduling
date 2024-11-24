package adapters

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"

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

func (r *PsqlTaskRepo) SaveOrUpdateTask(ctx context.Context, task model.Task) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		return q.UpsertTask(ctx, sqlc.UpsertTaskParams{
			ID:         uuid.UUID(task.Id()),
			Title:      task.Name(),
			Story:      task.Description(),
			ScheduleID: uuid.UUID(task.ScheduleId()),
		})
	})
}

func (r *PsqlTaskRepo) GetTask(ctx context.Context, id shared.Identity) (*model.Task, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	queries := sqlc.New(tx)

	dbTask, err := queries.GetTask(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}

	return model.NewTask(
		shared.Identity(dbTask.ID),
		shared.Identity(dbTask.ScheduleID),
		dbTask.Title,
		dbTask.Story,
	)
}

func (r *PsqlTaskRepo) DeleteTask(ctx context.Context, id shared.Identity) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		return q.DeleteTask(ctx, uuid.UUID(id))
	})
}
