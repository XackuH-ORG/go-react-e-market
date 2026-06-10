package main

import (
	"log/slog"
	"net/http"

	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
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
	r.Use(middleware.Timeout(app.config.HTTP.ReqTimeout))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("all good"))
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.HTTP.Addr,
		Handler:      h,
		WriteTimeout: app.config.HTTP.WriteTimeout,
		ReadTimeout:  app.config.HTTP.ReadTimeout,
		IdleTimeout:  app.config.HTTP.IdleTimeout,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks")
		// Здесь можно было бы ждать завершения фоновых задач (WaitGroups и т.д.)
		
		shutdownError <- nil
	}()

	app.logger.Info("Сервер запущен", "addr", app.config.HTTP.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("Сервер успешно остановлен")
	return nil
}

type application struct {
	config *config.Config
	logger *slog.Logger
	db     *pgxpool.Pool
}
