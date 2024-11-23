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

type PsqlWorkerRepo struct {
	db pgx.Conn
}

func (r PsqlWorkerRepo) Db() pgx.Conn {
	return r.db
}

func (r *PsqlWorkerRepo) SaveOrUpdateWorker(ctx context.Context, worker model.Worker) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		return q.UpsertWorker(ctx, sqlc.UpsertWorkerParams{
			ID:         uuid.UUID(worker.Id()),
			FirstName:  worker.FirstName(),
			LastName:   worker.LastName(),
			ScheduleID: uuid.UUID(worker.ScheudleId()),
		})
	})
}

func (r *PsqlWorkerRepo) GetWorker(ctx context.Context, id shared.Identity) (*model.Worker, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	queries := sqlc.New(tx)

	dbWorker, err := queries.GetWorker(ctx, uuid.UUID(id))

	if err != nil {
		return nil, err
	}

	return model.NewWorker(
		shared.Identity(dbWorker.ID),
		shared.Identity(dbWorker.ScheduleID),
		dbWorker.FirstName,
		dbWorker.LastName,
	)
}

func (r *PsqlWorkerRepo) DeleteWorker(ctx context.Context, id shared.Identity) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		return q.DeleteWorker(ctx, uuid.UUID(id))
	})
}
