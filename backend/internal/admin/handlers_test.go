package admin_test

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

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/admin"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/admin/mock"
	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

func setupAdminRouter(svc admin.AdminService) chi.Router {
	r := chi.NewRouter()
	h := admin.NewHandlers(svc)
	h.RegisterRoutes(r)
	return r
}

func TestHandlers_GetAllOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)

	t.Run("success", func(t *testing.T) {
		m.EXPECT().GetAllOrders(gomock.Any()).Return([]repo.Order{{ID: uuid.New()}}, nil)
		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("forbidden", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleCUSTOMER)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusForbidden, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().GetAllOrders(gomock.Any()).Return(nil, errors.New("db error"))
		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_SearchOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)

	t.Run("success", func(t *testing.T) {
		m.EXPECT().SearchOrders(gomock.Any(), "1234").Return([]repo.Order{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/orders/search?q=1234", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("bad request - invalid length", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders/search?q=123", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().SearchOrders(gomock.Any(), "1234").Return(nil, errors.New("db"))
		req := httptest.NewRequest(http.MethodGet, "/orders/search?q=1234", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_UpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)
	orderID := uuid.New()

	t.Run("success", func(t *testing.T) {
		m.EXPECT().UpdateOrderStatus(gomock.Any(), orderID, repo.OrderStatusDELIVERED).Return(nil)
		b, _ := json.Marshal(map[string]interface{}{"status": repo.OrderStatusDELIVERED})
		req := httptest.NewRequest(http.MethodPatch, "/orders/"+orderID.String()+"/status", bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/orders/invalid/status", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/orders/"+orderID.String()+"/status", bytes.NewReader([]byte("invalid")))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().UpdateOrderStatus(gomock.Any(), orderID, repo.OrderStatusDELIVERED).Return(errors.New("db"))
		b, _ := json.Marshal(map[string]interface{}{"status": repo.OrderStatusDELIVERED})
		req := httptest.NewRequest(http.MethodPatch, "/orders/"+orderID.String()+"/status", bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)

	t.Run("success", func(t *testing.T) {
		m.EXPECT().GetUsers(gomock.Any()).Return([]repo.User{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().GetUsers(gomock.Any()).Return(nil, errors.New("err"))
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_UpdateUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		m.EXPECT().UpdateUserRole(gomock.Any(), userID, repo.UserRoleADMIN).Return(nil)
		b, _ := json.Marshal(map[string]interface{}{"role": repo.UserRoleADMIN})
		req := httptest.NewRequest(http.MethodPatch, "/users/"+userID.String()+"/role", bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/users/invalid/role", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/users/"+userID.String()+"/role", bytes.NewReader([]byte("invalid")))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().UpdateUserRole(gomock.Any(), userID, repo.UserRoleADMIN).Return(errors.New("db"))
		b, _ := json.Marshal(map[string]interface{}{"role": repo.UserRoleADMIN})
		req := httptest.NewRequest(http.MethodPatch, "/users/"+userID.String()+"/role", bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)

	categoryID := uuid.New()
	desc := "desc"
	body := map[string]interface{}{
		"category_id": categoryID,
		"name":        "name",
		"description": desc,
	}

	t.Run("success", func(t *testing.T) {
		m.EXPECT().CreateProduct(gomock.Any(), categoryID, "name", &desc).Return(repo.Product{}, nil)
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid")))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().CreateProduct(gomock.Any(), categoryID, "name", &desc).Return(repo.Product{}, errors.New("db"))
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)

	productID := uuid.New()
	categoryID := uuid.New()
	desc := "desc"
	body := map[string]interface{}{
		"category_id": categoryID,
		"name":        "name",
		"description": desc,
	}

	t.Run("success", func(t *testing.T) {
		m.EXPECT().UpdateProduct(gomock.Any(), productID, categoryID, "name", &desc).Return(nil)
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/products/invalid", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), bytes.NewReader([]byte("invalid")))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().UpdateProduct(gomock.Any(), productID, categoryID, "name", &desc).Return(errors.New("db"))
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)

	productID := uuid.New()

	t.Run("success", func(t *testing.T) {
		m.EXPECT().DeleteProduct(gomock.Any(), productID).Return(nil)
		req := httptest.NewRequest(http.MethodDelete, "/products/"+productID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/invalid", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().DeleteProduct(gomock.Any(), productID).Return(errors.New("db"))
		req := httptest.NewRequest(http.MethodDelete, "/products/"+productID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_CreateSku(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)

	productID := uuid.New()
	body := map[string]interface{}{
		"product_id": productID,
		"sku_code":   "SKU-1",
		"attributes": json.RawMessage(`{"k":"v"}`),
		"price":      100,
		"stock":      10,
	}

	t.Run("success", func(t *testing.T) {
		m.EXPECT().CreateSku(gomock.Any(), productID, "SKU-1", []byte(`{"k":"v"}`), int32(100), int32(10)).Return(repo.Sku{}, nil)
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/skus", bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/skus", bytes.NewReader([]byte("invalid")))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().CreateSku(gomock.Any(), productID, "SKU-1", []byte(`{"k":"v"}`), int32(100), int32(10)).Return(repo.Sku{}, errors.New("db"))
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/skus", bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_UpdateSku(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)

	skuID := uuid.New()
	body := map[string]interface{}{
		"sku_code":   "SKU-1",
		"attributes": json.RawMessage(`{"k":"v"}`),
		"price":      100,
		"stock":      10,
	}

	t.Run("success", func(t *testing.T) {
		m.EXPECT().UpdateSku(gomock.Any(), skuID, "SKU-1", []byte(`{"k":"v"}`), int32(100), int32(10)).Return(nil)
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/skus/"+skuID.String(), bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/skus/invalid", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/skus/"+skuID.String(), bytes.NewReader([]byte("invalid")))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().UpdateSku(gomock.Any(), skuID, "SKU-1", []byte(`{"k":"v"}`), int32(100), int32(10)).Return(errors.New("db"))
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/skus/"+skuID.String(), bytes.NewReader(b))
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandlers_DeleteSku(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockAdminService(ctrl)
	r := setupAdminRouter(m)

	skuID := uuid.New()

	t.Run("success", func(t *testing.T) {
		m.EXPECT().DeleteSku(gomock.Any(), skuID).Return(nil)
		req := httptest.NewRequest(http.MethodDelete, "/skus/"+skuID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/skus/invalid", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("internal error", func(t *testing.T) {
		m.EXPECT().DeleteSku(gomock.Any(), skuID).Return(errors.New("db"))
		req := httptest.NewRequest(http.MethodDelete, "/skus/"+skuID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_role", string(repo.UserRoleADMIN)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
