package orders

import (
	"encoding/json"
	"net/http"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/lib/logger/sl"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateOrder(r.Context(), userID, req)
	if err != nil {
		slog.Error("failed to create order", sl.Err(err))
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(order)
}

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.service.GetOrders(r.Context(), userID)
	if err != nil {
		slog.Error("failed to get orders", sl.Err(err))
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(orders)
}

func (h *handler) GetAdminOrders(w http.ResponseWriter, r *http.Request) {
	searchIndex := r.URL.Query().Get("search")
	orders, err := h.service.GetAdminOrders(r.Context(), searchIndex)
	if err != nil {
		slog.Error("failed to get admin orders", sl.Err(err))
		http.Error(w, "Failed to get admin orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(orders)
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

func (h *handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	order, err := h.service.UpdateOrderStatus(r.Context(), orderID, req.Status)
	if err != nil {
		slog.Error("failed to update order status", sl.Err(err))
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(order)
}
