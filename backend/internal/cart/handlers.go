package cart

import (
	"encoding/json"
	"net/http"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/middleware"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/lib/logger/sl"
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

type AddToCartRequest struct {
	SkuID    uuid.UUID `json:"sku_id"`
	Quantity int32     `json:"quantity"`
}

func (h *handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	item, err := h.service.AddToCart(r.Context(), userID, req.SkuID, req.Quantity)
	if err != nil {
		slog.Error("failed to add to cart", sl.Err(err), slog.String("user_id", userID.String()))
		http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}

func (h *handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	items, err := h.service.GetCart(r.Context(), userID)
	if err != nil {
		slog.Error("failed to get cart", sl.Err(err), slog.String("user_id", userID.String()))
		http.Error(w, "Failed to get cart", http.StatusInternalServerError)
		return
	}

	if items == nil {
		items = make([]database.GetCartItemsRow, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}

func (h *handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	skuIDParam := chi.URLParam(r, "sku_id")
	skuID, err := uuid.Parse(skuIDParam)
	if err != nil {
		http.Error(w, "Invalid SKU ID", http.StatusBadRequest)
		return
	}

	err = h.service.RemoveFromCart(r.Context(), userID, skuID)
	if err != nil {
		slog.Error("failed to remove from cart", sl.Err(err), slog.String("user_id", userID.String()))
		http.Error(w, "Failed to remove from cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.service.ClearCart(r.Context(), userID)
	if err != nil {
		slog.Error("failed to clear cart", sl.Err(err), slog.String("user_id", userID.String()))
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
