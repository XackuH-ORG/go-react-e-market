package logger

import (
	"log/slog"
	"os"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/lib/logger/handlers/slogpretty"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

// SetupLogger инициализирует и возвращает объект логгера в зависимости от окружения.
func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		// PrettySlog
		log = setupPrettySlog()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		// По умолчанию используем настройки для прода
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
