package adapters

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"

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

func NewPsqlLocationRepo(db pgx.Conn) *PsqlLocationRepo {
	return &PsqlLocationRepo{db: db}
}

func (r *PsqlLocationRepo) SaveOrUpdateLocation(ctx context.Context, location model.Location) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		id := uuid.UUID(location.Id())
		scheduleId := uuid.UUID(location.ScheduleId())

		return q.UpsertLocation(ctx, sqlc.UpsertLocationParams{
			ID:         &id,
			Title:      location.Name(),
			Story:      location.Description(),
			ScheduleID: &scheduleId,
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

	locationId := uuid.UUID(id)
	dbLocation, err := queries.GetLocation(ctx, &locationId)

	if err != nil {
		return nil, err
	}

	return model.NewLocation(
		shared.Identity(*dbLocation.ID),
		shared.Identity(*dbLocation.ScheduleID),
		dbLocation.Title,
		dbLocation.Story,
	)
}

func (r *PsqlLocationRepo) DeleteLocation(ctx context.Context, id shared.Identity) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		locationId := uuid.UUID(id)
		return q.DeleteLocation(ctx, &locationId)
	})
}

var _ ports.LocationRepository = &PsqlLocationRepo{}
