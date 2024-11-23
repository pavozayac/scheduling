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

type PsqlLocationRepo struct {
	db pgx.Conn
}

func (r PsqlLocationRepo) Db() pgx.Conn {
	return r.db
}

func (r *PsqlLocationRepo) SaveOrUpdateLocation(ctx context.Context, location model.Location) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		return q.UpsertLocation(ctx, sqlc.UpsertLocationParams{
			ID:         uuid.UUID(location.Id()),
			Title:      location.Name(),
			Story:      location.Description(),
			ScheduleID: uuid.UUID(location.ScheduleId()),
		})
	})
}

func (r *PsqlLocationRepo) GetLocation(ctx context.Context, id shared.Identity) (*model.Location, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	queries := sqlc.New(tx)

	dbLocation, err := queries.GetLocation(ctx, uuid.UUID(id))

	if err != nil {
		return nil, err
	}

	return model.NewLocation(
		shared.Identity(dbLocation.ID),
		shared.Identity(dbLocation.ScheduleID),
		dbLocation.Title,
		dbLocation.Story,
	)
}

func (r *PsqlLocationRepo) DeleteLocation(ctx context.Context, id shared.Identity) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		return q.DeleteLocation(ctx, uuid.UUID(id))
	})
}
