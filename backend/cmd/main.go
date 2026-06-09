package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=admin dbname=emarket sslmode=disable"),
		},
	}

	// NOTE: logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// NOTE: База данных
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("Успешное подключение к базе данных", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
	}

	// NOTE: Запуск сервера
	if err := api.run(api.mount()); err != nil {
		slog.Error("Сервер не запустился!", "error", err)
		os.Exit(1)
	}
}
