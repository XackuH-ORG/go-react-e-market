package admin

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type Handlers struct {
	svc AdminService
}

func NewHandlers(svc AdminService) *Handlers {
	return &Handlers{svc: svc}
}

// AdminMiddleware checks if the user has ADMIN role.
// This is a temporary simple check using headers. Real auth middleware should set role in context.
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var role string
		ctxRole := r.Context().Value("user_role")
		if ctxRole != nil {
			if str, ok := ctxRole.(string); ok {
				role = str
			}
		}

		if role != string(repo.UserRoleADMIN) {
			http.Error(w, "forbidden: admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Use(AdminMiddleware)

	r.Get("/orders", h.GetAllOrders)
	r.Get("/orders/search", h.SearchOrders)
	r.Patch("/orders/{id}/status", h.UpdateOrderStatus)

	r.Get("/users", h.GetUsers)
	r.Patch("/users/{id}/role", h.UpdateUserRole)

	r.Post("/products", h.CreateProduct)
	r.Put("/products/{id}", h.UpdateProduct)
	r.Delete("/products/{id}", h.DeleteProduct)

	r.Post("/skus", h.CreateSku)
	r.Put("/skus/{id}", h.UpdateSku)
	r.Delete("/skus/{id}", h.DeleteSku)
}

func (h *Handlers) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.GetAllOrders(r.Context())
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

func (h *Handlers) SearchOrders(w http.ResponseWriter, r *http.Request) {
	index := r.URL.Query().Get("q")
	if len(index) != 4 {
		http.Error(w, "query must be exactly 4 characters", http.StatusBadRequest)
		return
	}

	orders, err := h.svc.SearchOrders(r.Context(), index)
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

func (h *Handlers) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		Status repo.OrderStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateOrderStatus(r.Context(), id, req.Status); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if users == nil {
		users = []repo.User{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handlers) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		Role repo.UserRole `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateUserRole(r.Context(), id, req.Role); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type CreateProductReq struct {
	CategoryID  uuid.UUID `json:"category_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
}

func (h *Handlers) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p, err := h.svc.CreateProduct(r.Context(), req.CategoryID, req.Name, req.Description)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *Handlers) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req CreateProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateProduct(r.Context(), id, req.CategoryID, req.Name, req.Description); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteProduct(r.Context(), id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type CreateSkuReq struct {
	ProductID  uuid.UUID       `json:"product_id"`
	SkuCode    string          `json:"sku_code"`
	Attributes json.RawMessage `json:"attributes"`
	Price      int32           `json:"price"`
	Stock      int32           `json:"stock"`
}

func (h *Handlers) CreateSku(w http.ResponseWriter, r *http.Request) {
	var req CreateSkuReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sku, err := h.svc.CreateSku(r.Context(), req.ProductID, req.SkuCode, req.Attributes, req.Price, req.Stock)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sku)
}

func (h *Handlers) UpdateSku(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req CreateSkuReq // We can reuse for put, just ignoring product_id mostly or accepting it.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateSku(r.Context(), id, req.SkuCode, req.Attributes, req.Price, req.Stock); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) DeleteSku(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteSku(r.Context(), id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
