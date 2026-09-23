package testutil

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(t *testing.T, ctx context.Context) *PostgresContainer {
	t.Helper()

	ctr, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join(getProjectRoot(), "init-db", "schema.sql")),
		postgres.WithDatabase("jobs"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("gojobs"),
		postgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, ctr)
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to retrieve connection string: %v", err)
	}

	return &PostgresContainer{
		PostgresContainer: ctr,
		ConnectionString:  connStr,
	}
}

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir // find project root
		}
		parent := filepath.Dir(dir)
		if parent == dir { // found o.s root and did not find project root - error
			break
		}
		dir = parent
	}
	return "."
}
