package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/repository"
)

// testDB - set one db for all tests
var testDB *sql.DB

var errAny = fmt.Errorf("any error")

// TestMain runs before all repository tests to:
// - create one instance for all tests
// - migrate up one time for all tests
// - kill and clean up instance
func TestMain(m *testing.M) {
	ctx := context.Background()

	// start container
	container, db, dsn, err := startPostgres(ctx)
	if err != nil {
		slog.Error("failed to start postgres container", "error", err)
		os.Exit(1)
	}

	defer db.Close()
	defer container.Terminate(ctx)

	testDB = db

	// migrates up
	if err := applyMigrations(dsn); err != nil {
		slog.Error("failed to apply migrations", "error", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// startPostgres sets up instance of postgres container and returns connection
func startPostgres(ctx context.Context) (testcontainers.Container, *sql.DB, string, error) {
	//testcontainers.GenericContainer creates and runs docker container for tests
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "testdb",
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
			},

			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("start container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, "", fmt.Errorf("got host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, "", fmt.Errorf("get port: %w", err)
	}

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, "", fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, "", fmt.Errorf("ping db: %w", err)
	}

	return container, db, dsn, nil
}

func applyMigrations(dsn string) error {
	m, err := migrate.New(
		"file://../../../migrations",
		dsn,
	)
	if err != nil {
		return fmt.Errorf("create migrations: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

// cleanDB cleans db tables between tests
func cleanDB(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(`
				TRUNCATE
					transactions,
					clients,
					businesses,
					business_bonus_settings
				RESTART IDENTITY CASCADE
	`)
	require.NoError(t, err)
}

func createTestClient(t *testing.T, phone, name string) domain.Client {
	t.Helper()
	client, err := repository.NewClientRepository(testDB).Create(
		context.Background(), phone, name,
	)
	require.NoError(t, err)
	return client
}

func createTestBusiness(t *testing.T, name, ownerPhone string) domain.Business {
	t.Helper()
	business, err := repository.NewBusinessRepository(testDB).Create(
		context.Background(), name, ownerPhone,
	)
	require.NoError(t, err)
	return business
}

func createTestBusinessWithBonusBalance(t *testing.T, name, ownerPhone string, delta int64) domain.Business {
	t.Helper()
	business := createTestBusiness(t, name, ownerPhone)
	updated, err := repository.NewBusinessRepository(testDB).UpdateBonusBalance(
		context.Background(), business.ID, delta,
	)
	require.NoError(t, err)
	assert.Equal(t, delta, updated.BonusBalance)
	return updated
}

func createTestTransaction(
	t *testing.T, 
	clientID, 
	businessID string, 
	amount int64, 
	bonusPercent float64, 
	description *string,
) domain.Transaction{
	t.Helper()
	tx, err := repository.NewTransactionRepository(testDB).CreateWithBonus(
		context.Background(), clientID, businessID, amount, bonusPercent, description,
	)
	require.NoError(t, err)
	return tx
}