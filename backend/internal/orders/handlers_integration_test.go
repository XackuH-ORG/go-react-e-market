//go:build integration

package orders_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/middleware"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/orders"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/testutil"
)

// Define a stub service for the test since we haven't implemented it yet.
// Wait, we are writing integration tests for handlers. But to test handlers, we need the real service and real DB.
// Since we are writing tests *first*, we will use the real service and expect it to fail if not implemented.

type testEnv struct {
	router *chi.Mux
	db     *repo.Queries
	userID uuid.UUID
	skuID  uuid.UUID
}

func setupOrderTestEnv(t *testing.T) testEnv {
	pool := testutil.SetupTestDB(t)
	db := repo.New(pool)
	ctx := context.Background()

	// 1. Create a test user
	user, err := db.CreateUser(ctx, repo.CreateUserParams{
		Email:        "order_test@example.com",
		PasswordHash: "hash",
		Role:         repo.UserRoleCUSTOMER,
	})
	require.NoError(t, err)

	// 2. Create product and sku
	product, err := db.CreateProduct(ctx, repo.CreateProductParams{
		Name:        "Order Test Product",
		Description: pgtype.Text{String: "For testing orders", Valid: true},
		ImageUrl:    pgtype.Text{String: "", Valid: false},
	})
	require.NoError(t, err)

	sku, err := db.CreateSku(ctx, repo.CreateSkuParams{
		ProductID: product.ID,
		SkuCode:   "ORDER-SKU-1",
		Price:     1000,
		Stock:     10,
	})
	require.NoError(t, err)

	// 3. Add to cart (orders will be created from cart)
	_, err = db.AddToCart(ctx, repo.AddToCartParams{
		UserID:   user.ID,
		SkuID:    sku.ID,
		Quantity: 2,
	})
	require.NoError(t, err)

	// 4. Setup Service & Handler
	// NOTE: We need a transaction-aware service here, so we will pass the pool.
	// For now we just pass nil to the service to see it fail, or we can pass a dummy.
	// Since we haven't implemented NewService with pool, let's just pass db.
	svc := orders.NewService(pool) // We will implement NewService to take *pgxpool.Pool
	h := orders.NewHandler(svc)

	r := chi.NewRouter()

	// Inject the test user ID into context for all routes
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.UserIDKey, user.ID)
			ctx = context.WithValue(ctx, middleware.UserRoleKey, user.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Post("/api/v1/orders", h.CreateOrder)
	r.Get("/api/v1/orders", h.GetOrders)

	// Admin routes
	r.Get("/api/v1/admin/orders", h.GetAdminOrders)
	r.Patch("/api/v1/admin/orders/{id}/status", h.UpdateStatus)

	return testEnv{
		router: r,
		db:     db,
		userID: user.ID,
		skuID:  sku.ID,
	}
}

func TestIntegration_OrderHandlers(t *testing.T) {
	env := setupOrderTestEnv(t)
	var createdOrderID uuid.UUID

	t.Run("Create Order - Success", func(t *testing.T) {
		reqBody := orders.CreateOrderRequest{
			DeliveryType:    "DELIVERY",
			DeliveryAddress: "123 Test St",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		env.router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusCreated, rr.Code)

		var resp repo.Order
		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Equal(t, "DELIVERY", string(resp.DeliveryType))
		assert.Equal(t, "123 Test St", resp.DeliveryAddress.String)
		assert.Equal(t, int32(2000), resp.TotalAmount) // 2 * 1000

		createdOrderID = resp.ID
	})

	t.Run("Get Orders", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
		rr := httptest.NewRecorder()

		env.router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp []repo.Order
		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)

		require.Len(t, resp, 1)
		assert.Equal(t, createdOrderID, resp[0].ID)
	})

	t.Run("Update Order Status (Admin)", func(t *testing.T) {
		reqBody := orders.UpdateStatusRequest{
			Status: "PROCESSING",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/orders/"+createdOrderID.String()+"/status", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		env.router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp repo.Order
		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, "PROCESSING", string(resp.Status))
	})
}
