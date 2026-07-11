package testutil

import (
	"context"
	"testing"
)

func TestSetupTestDB(t *testing.T) {
	// Данный тест загружает image и запускает docker-контейнер,
	// затем запускает миграции и проверяет всю работу
	pool := SetupTestDB(t)

	// Пинг БД
	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	t.Log("Successfully connected to testcontainers PostgreSQL 18")
}
