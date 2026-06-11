package orders_test

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
	"go.uber.org/mock/gomock"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/orders"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/orders/mock"
	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

func setupOrdersRouter(svc orders.OrdersService) chi.Router {
	r := chi.NewRouter()
	h := orders.NewHandlers(svc)
	h.RegisterRoutes(r)
	return r
}

func TestHandlers_Checkout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockOrdersService(ctrl)
	r := setupOrdersRouter(m)

	userID := uuid.New()
	addr := "123 Main St"
	expectedOrder := repo.Order{ID: uuid.New()}

	t.Run("success", func(t *testing.T) {
		m.EXPECT().Checkout(gomock.Any(), userID, repo.DeliveryTypeDELIVERY, &addr, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return(expectedOrder, nil)

		body := map[string]interface{}{
			"delivery_type":    repo.DeliveryTypeDELIVERY,
			"delivery_address": addr,
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewReader(b))
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var res repo.Order
		require.NoError(t, json.NewDecoder(w.Body).Decode(&res))
		assert.Equal(t, expectedOrder, res)
	})

	t.Run("unauthorized", func(t *testing.T) {
		body := map[string]interface{}{
			"delivery_type":    repo.DeliveryTypeDELIVERY,
			"delivery_address": addr,
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewReader(b))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().Checkout(gomock.Any(), userID, repo.DeliveryTypeDELIVERY, &addr, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return(repo.Order{}, errors.New("db error"))

		body := map[string]interface{}{
			"delivery_type":    repo.DeliveryTypeDELIVERY,
			"delivery_address": addr,
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewReader(b))
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockOrdersService(ctrl)
	r := setupOrdersRouter(m)

	userID := uuid.New()
	expectedOrders := []repo.Order{{ID: uuid.New()}}

	t.Run("success", func(t *testing.T) {
		m.EXPECT().GetUserOrders(gomock.Any(), userID).Return(expectedOrders, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var res []repo.Order
		require.NoError(t, json.NewDecoder(w.Body).Decode(&res))
		assert.Equal(t, expectedOrders, res)
	})

	t.Run("success - context user_id", func(t *testing.T) {
		m.EXPECT().GetUserOrders(gomock.Any(), userID).Return(expectedOrders, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := context.WithValue(req.Context(), "user_id", userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().GetUserOrders(gomock.Any(), userID).Return(nil, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", userID.String())
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
