package adapters

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/model"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/ports"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/domain/shared"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/infrastructure/db/sqlc"
	ishared "github.com/pavozayac/scheduling/src/constraint-service/internal/infrastructure/shared"
)

type PsqlScheduleRepo struct {
	db pgx.Conn
}

func (r PsqlScheduleRepo) Db() pgx.Conn {
	return r.db
}

func NewPsqlScheduleRepo(db pgx.Conn) *PsqlScheduleRepo {
	return &PsqlScheduleRepo{db: db}
}

func (r *PsqlScheduleRepo) SaveOrUpdateSchedule(ctx context.Context, schedule model.Schedule) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		// Save/update schedule
		scheduleID := uuid.UUID(schedule.Id())
		err := q.UpsertSchedule(ctx, sqlc.UpsertScheduleParams{
			ID:    &scheduleID,
			Title: schedule.Title(),
		})
		if err != nil {
			return err
		}

		// Handle constraints
		// First delete existing constraints
		err = q.DeleteScheduleConstraints(ctx, &scheduleID)
		if err != nil {
			return err
		}

		// Insert new constraints
		for _, constraint := range schedule.Constraints() {
			err = q.InsertConstraint(ctx, sqlc.InsertConstraintParams{
				ScheduleID: nillifyIdentity(constraint.ScheduleId()),
				LocationID: nillifyIdentity(constraint.LocationId()),
				TaskID:     nillifyIdentity(constraint.TaskId()),
				WorkerID:   nillifyIdentity(constraint.WorkerId()),
				StartSlot:  pgtype.Int4{Int32: int32(constraint.StartTime()), Valid: true},
				EndSlot:    pgtype.Int4{Int32: int32(constraint.EndTime()), Valid: true},
				Kind:       sqlc.ConstraintType(constraint.Type()),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PsqlScheduleRepo) GetSchedule(ctx context.Context, id shared.Identity) (*model.Schedule, error) {
	var schedule model.Schedule
	err := ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		// Get schedule
		scheduleID := uuid.UUID(id)
		sch, err := q.GetSchedule(ctx, &scheduleID)
		if err != nil {
			return err
		}

		// Get constraints
		constraints, err := q.GetAllConstraintsForSchedule(ctx, &scheduleID)
		if err != nil {
			return err
		}

		// Map DB constraints to domain constraints
		domainConstraints := make([]model.Constraint, 0)
		for _, c := range constraints {
			constraint, err := model.NewConstraint(
				shared.Identity(denillifyUuid(c.ScheduleID)),
				shared.Identity(denillifyUuid(c.WorkerID)),
				shared.Identity(denillifyUuid(c.TaskID)),
				shared.Identity(denillifyUuid(c.LocationID)),
				int(c.StartSlot.Int32),
				int(c.EndSlot.Int32),
				model.ConstraintType(c.Kind),
			)
			if err != nil {
				return err
			}
			domainConstraints = append(domainConstraints, constraint)
		}

		// Create domain schedule
		s, err := model.NewSchedule(shared.Identity(*sch.ID), sch.Title, domainConstraints)
		if err != nil {
			return err
		}

		schedule = *s

		return nil
	})

	return &schedule, err
}

func (r *PsqlScheduleRepo) DeleteSchedule(ctx context.Context, id shared.Identity) error {
	return ishared.TransactionDecorator(ctx, r, func(q *sqlc.Queries, ctx context.Context) error {
		// Delete constraints first due to foreign key
		scheduleID := uuid.UUID(id)
		err := q.DeleteScheduleConstraints(ctx, &scheduleID)
		if err != nil {
			return err
		}

		// Delete schedule
		return q.DeleteSchedule(ctx, &scheduleID)
	})
}

var _ ports.ScheduleRepository = &PsqlScheduleRepo{}
