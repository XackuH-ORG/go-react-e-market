package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/config"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/pkg/logger"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	// NOTE: Загрузка конфигурации (YAML + Env)
	cfg := config.MustLoad()

	// NOTE: Настройка логгера
	log := logger.SetupLogger(cfg.Env)
	slog.SetDefault(log)

	log.Info("Запуск приложения", "env", cfg.Env)

	// NOTE: База данных
	conn, err := pgx.Connect(ctx, cfg.DB.DSN)
	if err != nil {
		log.Error("Не удалось подключиться к базе данных", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	log.Info("Успешное подключение к базе данных")

	api := application{
		config: cfg,
		logger: log,
		db:     conn,
	}

	// NOTE: Запуск сервера
	if err := api.run(api.mount()); err != nil {
		log.Error("Сервер не запустился!", "error", err)
		os.Exit(1)
	}
}

