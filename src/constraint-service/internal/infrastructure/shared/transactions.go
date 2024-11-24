package shared

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/infrastructure/db/sqlc"
)

type PsqlDatabaser interface {
	Db() pgx.Conn
}

type QueryFunc func(*sqlc.Queries, context.Context) error

func TransactionDecorator(ctx context.Context, r PsqlDatabaser, queryFunc QueryFunc) error {
	db := r.Db()
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	queries := sqlc.New(tx)

	err = queryFunc(queries, ctx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}
