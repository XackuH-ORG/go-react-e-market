package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type mockCatalogService struct {
	GetCategoriesFunc func(ctx context.Context) ([]repo.Category, error)
	GetProductsFunc   func(ctx context.Context) ([]ProductResponse, error)
	GetProductFunc    func(ctx context.Context, id uuid.UUID) (*ProductResponse, error)
}

func (m *mockCatalogService) GetCategories(ctx context.Context) ([]repo.Category, error) {
	return m.GetCategoriesFunc(ctx)
}

func (m *mockCatalogService) GetProducts(ctx context.Context) ([]ProductResponse, error) {
	return m.GetProductsFunc(ctx)
}

func (m *mockCatalogService) GetProduct(ctx context.Context, id uuid.UUID) (*ProductResponse, error) {
	return m.GetProductFunc(ctx, id)
}

func TestHandlers_GetCategories(t *testing.T) {
	tests := []struct {
		name       string
		mockFunc   func(ctx context.Context) ([]repo.Category, error)
		wantStatus int
		wantBody   string
	}{
		{
			name: "Success",
			mockFunc: func(ctx context.Context) ([]repo.Category, error) {
				return []repo.Category{{ID: uuid.Nil, Name: "Cat1"}}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":"00000000-0000-0000-0000-000000000000","name":"Cat1","slug":"","parent_id":null,"created_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "Error",
			mockFunc: func(ctx context.Context) ([]repo.Category, error) {
				return nil, errors.New("error")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCatalogService{GetCategoriesFunc: tt.mockFunc}
			h := NewHandlers(svc)

			req := httptest.NewRequest("GET", "/categories", nil)
			w := httptest.NewRecorder()

			h.GetCategories(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			} else {
				assert.Equal(t, tt.wantBody, w.Body.String())
			}
		})
	}
}

func TestHandlers_GetProducts(t *testing.T) {
	tests := []struct {
		name       string
		mockFunc   func(ctx context.Context) ([]ProductResponse, error)
		wantStatus int
		wantBody   string
	}{
		{
			name: "Success",
			mockFunc: func(ctx context.Context) ([]ProductResponse, error) {
				return []ProductResponse{{ID: uuid.Nil, Name: "Prod1", Skus: []SkuResponse{}}}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":"00000000-0000-0000-0000-000000000000","category_id":"00000000-0000-0000-0000-000000000000","name":"Prod1","skus":[]}]`,
		},
		{
			name: "Error",
			mockFunc: func(ctx context.Context) ([]ProductResponse, error) {
				return nil, errors.New("error")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCatalogService{GetProductsFunc: tt.mockFunc}
			h := NewHandlers(svc)

			req := httptest.NewRequest("GET", "/products", nil)
			w := httptest.NewRecorder()

			h.GetProducts(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			} else {
				assert.Equal(t, tt.wantBody, w.Body.String())
			}
		})
	}
}

func TestHandlers_GetProduct(t *testing.T) {
	prodID := uuid.New()

	tests := []struct {
		name       string
		id         string
		mockFunc   func(ctx context.Context, id uuid.UUID) (*ProductResponse, error)
		wantStatus int
		wantBody   string
	}{
		{
			name: "Success",
			id:   prodID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*ProductResponse, error) {
				return &ProductResponse{ID: prodID, Name: "Prod1", Skus: []SkuResponse{}}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"` + prodID.String() + `","category_id":"00000000-0000-0000-0000-000000000000","name":"Prod1","skus":[]}`,
		},
		{
			name: "Not Found",
			id:   prodID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*ProductResponse, error) {
				return nil, nil
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "product not found\n",
		},
		{
			name:       "Invalid ID",
			id:         "invalid",
			mockFunc:   nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid product id\n",
		},
		{
			name: "Internal Error",
			id:   prodID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*ProductResponse, error) {
				return nil, errors.New("err")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCatalogService{GetProductFunc: tt.mockFunc}
			h := NewHandlers(svc)

			r := chi.NewRouter()
			h.RegisterRoutes(r)

			req := httptest.NewRequest("GET", "/products/"+tt.id, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var got map[string]interface{}
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				assert.Equal(t, prodID.String(), got["id"])
			} else {
				assert.Equal(t, tt.wantBody, w.Body.String())
			}
		})
	}
}
