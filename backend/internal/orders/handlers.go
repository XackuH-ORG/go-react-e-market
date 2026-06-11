package orders

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type Handlers struct {
	svc OrdersService
}

func NewHandlers(svc OrdersService) *Handlers {
	return &Handlers{svc: svc}
}

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Post("/checkout", h.Checkout)
	r.Get("/", h.GetUserOrders)
}

func getUserID(r *http.Request) (uuid.UUID, error) {
	headerID := r.Header.Get("X-User-ID")
	if headerID != "" {
		return uuid.Parse(headerID)
	}
	val := r.Context().Value("user_id")
	if val != nil {
		if str, ok := val.(string); ok {
			return uuid.Parse(str)
		}
	}
	return uuid.Nil, errors.New("unauthorized: missing user_id")
}

type CheckoutRequest struct {
	DeliveryType    repo.DeliveryType `json:"delivery_type"`
	DeliveryAddress *string           `json:"delivery_address,omitempty"`
	PickupPointID   *uuid.UUID        `json:"pickup_point_id,omitempty"`
	PromoCodeID     *uuid.UUID        `json:"promo_code_id,omitempty"`
}

func (h *Handlers) Checkout(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.svc.Checkout(r.Context(), userID, req.DeliveryType, req.DeliveryAddress, req.PickupPointID, req.PromoCodeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *Handlers) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	orders, err := h.svc.GetUserOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if orders == nil {
		orders = []repo.Order{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
