package adapters

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pavozayac/constraints/internal/domain/model"
	"github.com/pavozayac/constraints/internal/domain/ports"
	"github.com/pavozayac/constraints/internal/domain/shared"
	"github.com/pavozayac/constraints/internal/infrastructure/db/sqlc"
	ishared "github.com/pavozayac/constraints/internal/infrastructure/shared"
)

type PsqlWorkerRepo struct {
	db pgx.Conn
}

func (r PsqlWorkerRepo) Db() pgx.Conn {
	return r.db
}

func NewPsqlWorkerRepo(db pgx.Conn) *PsqlWorkerRepo {
	return &PsqlWorkerRepo{db: db}
}

func (r *PsqlWorkerRepo) SaveOrUpdateWorker(ctx context.Context, worker model.Worker) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		id := uuid.UUID(worker.Id())
		scheduleID := uuid.UUID(worker.ScheduleId())
		return q.UpsertWorker(ctx, sqlc.UpsertWorkerParams{
			ID:         &id,
			FirstName:  worker.FirstName(),
			LastName:   worker.LastName(),
			ScheduleID: &scheduleID,
		})
	})
}

func (r *PsqlWorkerRepo) GetWorker(ctx context.Context, id shared.Identity) (worker *model.Worker, err error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = errors.Join(err, tx.Rollback(ctx))
	}()

	queries := sqlc.New(tx)

	uuidID := uuid.UUID(id)
	dbWorker, err := queries.GetWorker(ctx, &uuidID)

	if err != nil {
		return nil, err
	}

	return model.NewWorker(
		shared.Identity(*dbWorker.ID),
		shared.Identity(*dbWorker.ScheduleID),
		dbWorker.FirstName,
		dbWorker.LastName,
	)
}

func (r *PsqlWorkerRepo) DeleteWorker(ctx context.Context, id shared.Identity) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		uuidID := uuid.UUID(id)
		return q.DeleteWorker(ctx, &uuidID)
	})
}

var _ ports.WorkerRepository = &PsqlWorkerRepo{}
