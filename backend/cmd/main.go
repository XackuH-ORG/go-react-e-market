package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/config"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/lib/logger/sl"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// @title           Go-React E-Market API
// @version         1.0
// @description     API-сервер для интернет-магазина Go-React E-Market.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8082
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Введите "Bearer " и затем JWT токен в поле ввода.
func main() {
	ctx := context.Background()

	// NOTE: Загрузка конфигурации
	cfg := config.MustLoad()

	// NOTE: Настройка логгера
	log := logger.SetupLogger(cfg.Env)
	slog.SetDefault(log)

	log.Info("Запуск приложения", "env", cfg.Env)

	// NOTE: База данных
	conn, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		log.Error("Не удалось подключиться к базе данных", sl.Err(err))
		os.Exit(1)
	}
	defer conn.Close()

	log.Info("Успешное подключение к базе данных")

	api := application{
		config: cfg,
		logger: log,
		db:     conn,
	}

	// NOTE: Запуск сервера
	if err := api.run(api.mount()); err != nil {
		log.Error("Сервер не запустился!", sl.Err(err))
		os.Exit(1)
	}
}
