package testutil

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // для goose
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetupTestDB запускает PostgreSQL 18 container, запускает миграции goose, и возвращает подключение pgxpool.
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	dbName := "testdb"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate postgres container: %v", err)
		}
	})

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// 1. Запуск миграций
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("failed to open sql db for migrations: %v", err)
	}

	// Находим каталог миграций относительно этого файла.
	// Текущий файл: backend/internal/testutil/db.go
	// Миграции: backend/internal/adapters/postgresql/migrations
	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "..", "adapters", "postgresql", "migrations")

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("failed to set goose dialect: %v", err)
	}
	if err := goose.Up(db, migrationsDir); err != nil {
		t.Fatalf("failed to run goose migrations: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("failed to close sql db after migrations: %v", err)
	}

	// 2. Возвращаем pgxpool
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		t.Fatalf("failed to parse pgxpool config: %v", err)
	}
	poolConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		t.Fatalf("failed to connect pgxpool: %v", err)
	}

	// Закрываем пул после очистки теста
	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}
