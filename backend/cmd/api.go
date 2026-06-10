package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/config"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)              // важно для ограничения скорости запросов
	r.Use(middleware.ClientIPFromRemoteAddr) // важно для ограничения скорости, аналитики и отслеживания
	r.Use(middleware.Logger)                 // для логирования запросов, может быть отключен в продакшене
	r.Use(middleware.Recoverer)              // для восстановления после сбоев, может быть отключен в продакшене

	// Установка таймаута для контекста запроса (ctx), который сигнализирует
	// через ctx.Done() об истечении времени и требует остановки
	// дальнейшей обработки запроса.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	// Инициализируем репозиторий sqlc и прокидываем в сервис продуктов
	queries := repo.New(app.db)
	productService := products.NewService(queries)
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.HTTP.Addr,
		Handler:      h,
		WriteTimeout: app.config.HTTP.Timeout,
		ReadTimeout:  app.config.HTTP.Timeout * 2, // или можно задать явно, используем таймаут как базу
		IdleTimeout:  app.config.HTTP.IdleTimeout,
	}

	app.logger.Info("Сервер запущен", "addr", app.config.HTTP.Addr)

	return srv.ListenAndServe()
}

type application struct {
	config *config.Config
	logger *slog.Logger
	db     *pgx.Conn
}

