package shared

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupMigratedPostgresContainer(t *testing.T, ctx context.Context, migrationsDir string) (string, postgres.PostgresContainer, func(*testing.T)) {
	postgresContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithUsername("dbuser"),
		postgres.WithPassword("password"),
		postgres.WithDatabase("scheduling_db"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	require.NoError(t, err)

	connStr, err := postgresContainer.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := goose.OpenDBWithDriver("pgx", connStr)
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, goose.SetDialect("postgres"))
	require.NoError(t, goose.Up(db, migrationsDir))

	return connStr, *postgresContainer, func(t *testing.T) {
		testcontainers.CleanupContainer(t, postgresContainer)
		require.NoError(t, err)
	}
}
