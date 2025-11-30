package testdb

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	sharedContainer *PostgresContainer
	sharedOnce      sync.Once
	sharedMu        sync.Mutex
)

// SetupSharedPostgres creates a single PostgreSQL container shared across all tests
// This is much faster than creating a new container for each test.
// Use this for local development where speed > isolation.
//
// IMPORTANT: Tests using shared container CANNOT run in parallel!
//
// Usage:
//
//	func TestMyService(t *testing.T) {
//	    pgContainer := testdb.SetupSharedPostgres(t)
//	    defer pgContainer.Cleanup(t)  // ← Only call once at top level
//
//	    t.Run("Test1", func(t *testing.T) {
//	        testdb.CleanupTables(t, pgContainer.DB, "my_table")
//	        // ... test
//	    })
//	}
func SetupSharedPostgres(t *testing.T) *PostgresContainer {
	t.Helper()

	sharedOnce.Do(func() {
		ctx := context.Background()
		pgContainer, err := postgres.Run(ctx,
			"postgres:16-alpine",
			postgres.WithDatabase("testdb"),
			postgres.WithUsername("postgres"),
			postgres.WithPassword("postgres"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2),
			),
		)
		require.NoError(t, err)

		connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
		require.NoError(t, err)

		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))
		db := bun.NewDB(sqldb, pgdialect.New())

		err = db.Ping()
		require.NoError(t, err)

		sharedContainer = &PostgresContainer{
			Container: pgContainer,
			DB:        db,
			DSN:       connStr,
		}
	})

	return sharedContainer
}

func CleanupTables(t *testing.T, db *bun.DB, tables ...string) {
	t.Helper()

	ctx := context.Background()

	for _, table := range tables {
		_, err := db.ExecContext(ctx, "TRUNCATE "+table+" RESTART IDENTITY CASCADE")
		require.NoError(t, err, "failed to truncate table: %s", table)
	}
}
