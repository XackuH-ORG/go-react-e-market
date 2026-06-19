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

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/admin"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/auth"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/cart"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/catalog"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/orders"
	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)              // важно для ограничения скорости запросов
	r.Use(middleware.ClientIPFromRemoteAddr) // важно для ограничения скорости, аналитики и отслеживания
	r.Use(middleware.Logger)                 // для логирования запросов, может быть отключен в продакшене
	r.Use(middleware.Recoverer)              // для восстановления после сбоев, может быть отключен в продакшене

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // In production, replace with specific origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Security Headers Middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			next.ServeHTTP(w, r)
		})
	})

	// Установка таймаута для контекста запроса (ctx), который сигнализирует
	// через ctx.Done() об истечении времени и требует остановки
	// дальнейшей обработки запроса.
	r.Use(middleware.Timeout(app.config.HTTP.ReqTimeout))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("all good"))
	})

	q := repo.New(app.db)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("my-secret-key")
	}

	authSvc := auth.NewAuthService(q, jwtSecret)
	authHandlers := auth.NewHandlers(authSvc, jwtSecret)

	catalogSvc := catalog.NewCatalogService(q)
	catalogHandlers := catalog.NewHandlers(catalogSvc)

	cartSvc := cart.NewCartService(q)
	cartHandlers := cart.NewHandlers(cartSvc)

	ordersSvc := orders.NewOrdersService(app.db)
	ordersHandlers := orders.NewHandlers(ordersSvc)

	adminSvc := admin.NewAdminService(app.db, q)
	adminHandlers := admin.NewHandlers(adminSvc)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			authHandlers.RegisterRoutes(r)
		})
		// Catalog routes register their own /products and /categories prefixes
		catalogHandlers.RegisterRoutes(r)

		r.Route("/cart", func(r chi.Router) {
			cartHandlers.RegisterRoutes(r)
		})
		r.Route("/orders", func(r chi.Router) {
			ordersHandlers.RegisterRoutes(r)
		})
		r.Route("/admin", func(r chi.Router) {
			adminHandlers.RegisterRoutes(r)
		})
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
