//go:build integration

package products_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/products"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/testutil"
)

func setupTestRouter(t *testing.T) (*chi.Mux, repo.Querier) {
	pool := testutil.SetupTestDB(t)
	db := repo.New(pool)

	svc := products.NewService(db)
	h := products.NewHandler(svc)

	r := chi.NewRouter()
	r.Post("/api/v1/admin/products", h.CreateProduct)
	r.Post("/api/v1/admin/skus", h.CreateSku)
	r.Get("/api/v1/products", h.ListProducts)
	r.Get("/api/v1/products/{id}", h.GetProduct)
	r.Put("/api/v1/admin/products/{id}", h.UpdateProduct)
	r.Delete("/api/v1/admin/products/{id}", h.DeleteProduct)

	return r, db
}

func TestIntegration_ProductHandlers(t *testing.T) {
	router, _ := setupTestRouter(t)

	var createdProductID uuid.UUID

	t.Run("Create Product - Success", func(t *testing.T) {
		reqBody := products.CreateProductRequest{
			Name:        "Integration Test Product",
			Description: "Desc",
			ImageURL:    "http://test.com/img.png",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusCreated, rr.Code)

		var resp repo.Product
		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, reqBody.Name, resp.Name)
		assert.Equal(t, reqBody.Description, resp.Description.String)
		assert.Equal(t, reqBody.ImageURL, resp.ImageUrl.String)

		createdProductID = resp.ID
	})

	t.Run("Create Product - Invalid Input", func(t *testing.T) {
		reqBody := products.CreateProductRequest{
			Description: "No name",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Create SKU - Success", func(t *testing.T) {
		reqBody := products.CreateSkuRequest{
			ProductID: createdProductID,
			SkuCode:   "SKU-INT-1",
			Price:     1000,
			Stock:     50,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/skus", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusCreated, rr.Code)

		var resp repo.Sku
		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, reqBody.SkuCode, resp.SkuCode)
		assert.Equal(t, reqBody.Price, resp.Price)
		assert.Equal(t, reqBody.Stock, resp.Stock)
		assert.Equal(t, createdProductID, resp.ProductID)
	})

	t.Run("Get Product Details", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+createdProductID.String(), nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp products.ProductDetails
		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Equal(t, createdProductID, resp.Product.ID)
		assert.Equal(t, "Integration Test Product", resp.Product.Name)
		require.Len(t, resp.Skus, 1)
		assert.Equal(t, "SKU-INT-1", resp.Skus[0].SkuCode)
	})

	t.Run("List Products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products?limit=10", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp []repo.Product
		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)

		assert.NotEmpty(t, resp)

		found := false
		for _, p := range resp {
			if p.ID == createdProductID {
				found = true
				break
			}
		}
		assert.True(t, found, "Product should be in list")
	})

	t.Run("Update Product", func(t *testing.T) {
		reqBody := products.UpdateProductRequest{
			Name:        "Integration Test Product v2",
			Description: "Updated desc",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/products/"+createdProductID.String(), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp repo.Product
		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, "Integration Test Product v2", resp.Name)
		assert.Equal(t, "Updated desc", resp.Description.String)
	})

	t.Run("Delete Product", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/products/"+createdProductID.String(), nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNoContent, rr.Code)

		// Verify it is gone
		reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+createdProductID.String(), nil)
		rrGet := httptest.NewRecorder()
		router.ServeHTTP(rrGet, reqGet)

		require.Equal(t, http.StatusNotFound, rrGet.Code)
	})
}
