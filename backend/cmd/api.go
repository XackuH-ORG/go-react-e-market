package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	//TODO: добавить логгер в приложение и использовать его для логирования, вместо slog.Info
	slog.Info(fmt.Sprintf("Сервер запущен: %s", app.config.addr))

	return srv.ListenAndServe()
}

type application struct {
	config config
	// logger
	// db driver
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
