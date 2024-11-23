package shared

import (
	"context"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupMigratedPostgresContainer(t *testing.T, ctx context.Context) (postgres.PostgresContainer, func(*testing.T)) {
	postgresContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithUsername("dbuser"),
		postgres.WithPassword("password"),
		postgres.WithDatabase("scheduling_db"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	assert.NoError(t, err)

	connStr, err := postgresContainer.ConnectionString(ctx)
	assert.NoError(t, err)

	db, err := goose.OpenDBWithDriver("pgx", connStr)
	assert.NoError(t, err)
	defer db.Close()

	assert.NoError(t, err)
	goose.SetBaseFS(os.DirFS("../db"))

	assert.NoError(t, goose.SetDialect("postgres"))
	assert.NoError(t, goose.Up(db, "migrations"))

	return *postgresContainer, func(t *testing.T) {
		testcontainers.CleanupContainer(t, postgresContainer)
		assert.NoError(t, err)
	}
}
