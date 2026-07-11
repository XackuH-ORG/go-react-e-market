//go:build integration

package cart_test

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
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/cart"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/middleware"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/testutil"
)

func setupCartTestEnv(t *testing.T) (*chi.Mux, *repo.Queries, uuid.UUID, uuid.UUID) {
	pool := testutil.SetupTestDB(t)
	db := repo.New(pool)

	ctx := context.Background()

	// 1. Create a test user
	user, err := db.CreateUser(ctx, repo.CreateUserParams{
		Email:        "cart_test@example.com",
		PasswordHash: "hash",
		Role:         repo.UserRoleCUSTOMER,
	})
	require.NoError(t, err)

	// 2. Create a test product and SKU
	product, err := db.CreateProduct(ctx, repo.CreateProductParams{
		Name:        "Cart Test Product",
		Description: pgtype.Text{String: "For testing cart", Valid: true},
		ImageUrl:    pgtype.Text{String: "", Valid: false},
	})
	require.NoError(t, err)

	sku, err := db.CreateSku(ctx, repo.CreateSkuParams{
		ProductID: product.ID,
		SkuCode:   "CART-SKU-1",
		Price:     500,
		Stock:     100,
	})
	require.NoError(t, err)

	// 3. Setup Cart Handler
	svc := cart.NewService(db)
	h := cart.NewHandler(svc)

	r := chi.NewRouter()

	// Inject the test user ID into context for all routes
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.UserIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Post("/api/v1/cart", h.AddToCart)
	r.Get("/api/v1/cart", h.GetCart)
	r.Delete("/api/v1/cart/{sku_id}", h.RemoveFromCart)
	r.Delete("/api/v1/cart", h.ClearCart)

	return r, db, user.ID, sku.ID
}

func TestIntegration_CartHandlers(t *testing.T) {
	router, _, _, skuID := setupCartTestEnv(t)

	t.Run("Get Empty Cart", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		var items []repo.GetCartItemsRow
		err := json.NewDecoder(rr.Body).Decode(&items)
		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("Add to Cart - Success", func(t *testing.T) {
		reqBody := cart.AddToCartRequest{
			SkuID:    skuID,
			Quantity: 2,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/cart", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var item repo.CartItem
		err := json.NewDecoder(rr.Body).Decode(&item)
		require.NoError(t, err)
		assert.Equal(t, skuID, item.SkuID)
		assert.Equal(t, int32(2), item.Quantity)
	})

	t.Run("Add to Cart - Increment Quantity", func(t *testing.T) {
		reqBody := cart.AddToCartRequest{
			SkuID:    skuID,
			Quantity: 3,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/cart", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var item repo.CartItem
		err := json.NewDecoder(rr.Body).Decode(&item)
		require.NoError(t, err)
		assert.Equal(t, skuID, item.SkuID)
		assert.Equal(t, int32(5), item.Quantity) // 2 + 3 = 5
	})

	t.Run("Get Cart with Items", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var items []repo.GetCartItemsRow
		err := json.NewDecoder(rr.Body).Decode(&items)
		require.NoError(t, err)

		require.Len(t, items, 1)
		assert.Equal(t, skuID, items[0].SkuID)
		assert.Equal(t, int32(5), items[0].Quantity)
		assert.Equal(t, "Cart Test Product", items[0].Name)
	})

	t.Run("Remove from Cart", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/cart/"+skuID.String(), nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNoContent, rr.Code)

		// Verify it's removed
		reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
		rrGet := httptest.NewRecorder()
		router.ServeHTTP(rrGet, reqGet)

		var items []repo.GetCartItemsRow
		_ = json.NewDecoder(rrGet.Body).Decode(&items)
		assert.Empty(t, items)
	})

	t.Run("Clear Cart", func(t *testing.T) {
		// First, add item back
		reqBody := cart.AddToCartRequest{
			SkuID:    skuID,
			Quantity: 1,
		}
		bodyBytes, _ := json.Marshal(reqBody)
		reqPost := httptest.NewRequest(http.MethodPost, "/api/v1/cart", bytes.NewBuffer(bodyBytes))
		router.ServeHTTP(httptest.NewRecorder(), reqPost)

		// Then clear it
		reqDelete := httptest.NewRequest(http.MethodDelete, "/api/v1/cart", nil)
		rrDelete := httptest.NewRecorder()
		router.ServeHTTP(rrDelete, reqDelete)

		require.Equal(t, http.StatusNoContent, rrDelete.Code)

		// Verify empty
		reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
		rrGet := httptest.NewRecorder()
		router.ServeHTTP(rrGet, reqGet)

		var items []repo.GetCartItemsRow
		_ = json.NewDecoder(rrGet.Body).Decode(&items)
		assert.Empty(t, items)
	})
}
