package cart

import (
	"bytes"
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

type mockCartService struct {
	AddToCartFunc func(ctx context.Context, userID, skuID uuid.UUID, quantity int32) error
	GetCartFunc   func(ctx context.Context, userID uuid.UUID) ([]repo.GetCartItemsRow, error)
	ClearCartFunc func(ctx context.Context, userID uuid.UUID) error
}

func (m *mockCartService) AddToCart(ctx context.Context, userID, skuID uuid.UUID, quantity int32) error {
	return m.AddToCartFunc(ctx, userID, skuID, quantity)
}

func (m *mockCartService) GetCart(ctx context.Context, userID uuid.UUID) ([]repo.GetCartItemsRow, error) {
	return m.GetCartFunc(ctx, userID)
}

func (m *mockCartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	return m.ClearCartFunc(ctx, userID)
}

func TestHandlers_GetCart(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name       string
		headers    map[string]string
		mockFunc   func(ctx context.Context, userID uuid.UUID) ([]repo.GetCartItemsRow, error)
		wantStatus int
		wantBody   string
	}{
		{
			name:    "Success",
			headers: map[string]string{"X-User-ID": userID.String()},
			mockFunc: func(ctx context.Context, u uuid.UUID) ([]repo.GetCartItemsRow, error) {
				assert.Equal(t, userID, u)
				return []repo.GetCartItemsRow{{SkuID: uuid.Nil, Quantity: 2}}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   `[{"quantity":2,"sku_id":"00000000-0000-0000-0000-000000000000","price":0,"stock":0,"name":""}]`,
		},
		{
			name:       "Missing User ID",
			headers:    nil,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized: missing user_id\n",
		},
		{
			name:    "Error",
			headers: map[string]string{"X-User-ID": userID.String()},
			mockFunc: func(ctx context.Context, u uuid.UUID) ([]repo.GetCartItemsRow, error) {
				return nil, errors.New("error")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal error\n",
		},
		{
			name:    "Null Items Returns Empty Array",
			headers: map[string]string{"X-User-ID": userID.String()},
			mockFunc: func(ctx context.Context, u uuid.UUID) ([]repo.GetCartItemsRow, error) {
				return nil, nil // Return nil slice
			},
			wantStatus: http.StatusOK,
			wantBody:   `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCartService{GetCartFunc: tt.mockFunc}
			h := NewHandlers(svc)

			req := httptest.NewRequest("GET", "/", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()

			h.GetCart(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			} else {
				assert.Equal(t, tt.wantBody, w.Body.String())
			}
		})
	}
}

func TestHandlers_AddToCart(t *testing.T) {
	userID := uuid.New()
	skuID := uuid.New()

	tests := []struct {
		name       string
		headers    map[string]string
		reqBody    AddToCartRequest
		mockFunc   func(ctx context.Context, u, s uuid.UUID, q int32) error
		wantStatus int
	}{
		{
			name:    "Success",
			headers: map[string]string{"X-User-ID": userID.String()},
			reqBody: AddToCartRequest{SkuID: skuID, Quantity: 2},
			mockFunc: func(ctx context.Context, u, s uuid.UUID, q int32) error {
				assert.Equal(t, userID, u)
				assert.Equal(t, skuID, s)
				assert.Equal(t, int32(2), q)
				return nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Missing User ID",
			headers:    nil,
			reqBody:    AddToCartRequest{SkuID: skuID, Quantity: 2},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid Quantity",
			headers:    map[string]string{"X-User-ID": userID.String()},
			reqBody:    AddToCartRequest{SkuID: skuID, Quantity: 0},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "Error from Service",
			headers: map[string]string{"X-User-ID": userID.String()},
			reqBody: AddToCartRequest{SkuID: skuID, Quantity: 2},
			mockFunc: func(ctx context.Context, u, s uuid.UUID, q int32) error {
				return errors.New("err")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCartService{AddToCartFunc: tt.mockFunc}
			h := NewHandlers(svc)

			body, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/add", bytes.NewReader(body))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()

			h.AddToCart(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestHandlers_ClearCart(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name       string
		headers    map[string]string
		mockFunc   func(ctx context.Context, u uuid.UUID) error
		wantStatus int
	}{
		{
			name:    "Success",
			headers: map[string]string{"X-User-ID": userID.String()},
			mockFunc: func(ctx context.Context, u uuid.UUID) error {
				assert.Equal(t, userID, u)
				return nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Missing User ID",
			headers:    nil,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:    "Error from Service",
			headers: map[string]string{"X-User-ID": userID.String()},
			mockFunc: func(ctx context.Context, u uuid.UUID) error {
				return errors.New("err")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCartService{ClearCartFunc: tt.mockFunc}
			h := NewHandlers(svc)

			req := httptest.NewRequest("DELETE", "/", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()

			h.ClearCart(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestHandlers_RegisterRoutes(t *testing.T) {
	h := NewHandlers(&mockCartService{})
	r := chi.NewRouter()
	h.RegisterRoutes(r)
	
	assert.NotNil(t, r)
}
