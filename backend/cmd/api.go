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

	_ "github.com/XackuH-ORG/go-react-e-market/backend/docs"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/auth"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/cart"
	appMiddleware "github.com/XackuH-ORG/go-react-e-market/backend/internal/middleware"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/orders"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/products"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/users"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Для разработки разрешаем всё
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)              // важно для ограничения скорости запросов
	r.Use(middleware.ClientIPFromRemoteAddr) // важно для ограничения скорости, аналитики и отслеживания
	r.Use(middleware.Logger)                 // для логирования запросов, может быть отключен в продакшене
	r.Use(middleware.Recoverer)              // для восстановления после сбоев, может быть отключен в продакшене

	// Установка таймаута для контекста запроса (ctx), который сигнализирует
	// через ctx.Done() об истечении времени и требует остановки
	// дальнейшей обработки запроса.
	r.Use(middleware.Timeout(app.config.HTTP.ReqTimeout))

	r.Get("/health", app.health)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	// Инициализация сервисов и хендлеров
	queries := database.New(app.db)
	authService := auth.NewService(queries)
	authHandler := auth.NewHandler(authService)
	authMw := appMiddleware.NewAuthMiddleware(queries)

	usersService := users.NewService(queries)
	usersHandler := users.NewHandler(usersService)

	productsService := products.NewService(queries)
	productsHandler := products.NewHandler(productsService)

	cartService := cart.NewService(queries)
	cartHandler := cart.NewHandler(cartService)

	orderService := orders.NewService(app.db)
	orderHandler := orders.NewHandler(orderService)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		r.Get("/products", productsHandler.ListProducts)
		r.Get("/products/{id}", productsHandler.GetProduct)

		// Защищенные роуты для ВСЕХ авторизованных пользователей (Корзина, Заказы)
		r.Group(func(r chi.Router) {
			r.Use(authMw.RequireAuth)

			r.Post("/cart", cartHandler.AddToCart)
			r.Get("/cart", cartHandler.GetCart)
			r.Delete("/cart/{sku_id}", cartHandler.RemoveFromCart)
			r.Delete("/cart", cartHandler.ClearCart)

			r.Post("/orders", orderHandler.CreateOrder)
			r.Get("/orders", orderHandler.GetOrders)
		})

		// Защищенные роуты ТОЛЬКО ДЛЯ АДМИНОВ (Товары, Управление заказами)
		r.Group(func(r chi.Router) {
			r.Use(authMw.RequireAuth)
			r.Use(appMiddleware.RequireAdmin) // Второй слой защиты

			r.Get("/admin/users", usersHandler.GetUsers)
			r.Patch("/admin/users/{id}/role", usersHandler.UpdateRole)

			r.Post("/admin/products", productsHandler.CreateProduct)
			r.Put("/admin/products/{id}", productsHandler.UpdateProduct)
			r.Delete("/admin/products/{id}", productsHandler.DeleteProduct)
			r.Post("/admin/products/{id}/image", productsHandler.UploadImage)

			r.Post("/admin/skus", productsHandler.CreateSku)

			r.Get("/admin/orders", orderHandler.GetAdminOrders)
			r.Patch("/admin/orders/{id}/status", orderHandler.UpdateStatus)
		})
	})

	return r
}

// health godoc
// @Summary      Проверка здоровья API
// @Description  Возвращает статус работоспособности сервера бэкенда
// @Tags         system
// @Produce      plain
// @Success      200  {string}  string  "all good"
// @Router       /health [get]
func (app *application) health(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("all good"))
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
