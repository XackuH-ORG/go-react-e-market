package cart

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type Handlers struct {
	svc CartService
}

func NewHandlers(svc CartService) *Handlers {
	return &Handlers{svc: svc}
}

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Get("/", h.GetCart)
	r.Post("/add", h.AddToCart)
	r.Delete("/", h.ClearCart)
}

// getUserID is a helper to get user ID from context or header
func getUserID(r *http.Request) (uuid.UUID, error) {
	// For now, allow passing user ID via header if auth middleware is not fully hooked up
	headerID := r.Header.Get("X-User-ID")
	if headerID != "" {
		return uuid.Parse(headerID)
	}

	// Try from context
	val := r.Context().Value("user_id")
	if val != nil {
		if str, ok := val.(string); ok {
			return uuid.Parse(str)
		}
	}

	return uuid.Nil, errors.New("unauthorized: missing user_id")
}

type AddToCartRequest struct {
	SkuID    uuid.UUID `json:"sku_id"`
	Quantity int32     `json:"quantity"`
}

func (h *Handlers) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "quantity must be greater than zero", http.StatusBadRequest)
		return
	}

	if err := h.svc.AddToCart(r.Context(), userID, req.SkuID, req.Quantity); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	items, err := h.svc.GetCart(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// If no items, return empty array instead of null
	if items == nil {
		items = []repo.GetCartItemsRow{} // wait repo is not imported here, let me fix it
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *Handlers) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := h.svc.ClearCart(r.Context(), userID); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
