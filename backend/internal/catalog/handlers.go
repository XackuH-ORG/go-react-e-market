package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handlers struct {
	svc CatalogService
}

func NewHandlers(svc CatalogService) *Handlers {
	return &Handlers{svc: svc}
}

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Get("/categories", h.GetCategories)
	r.Get("/products", h.GetProducts)
	r.Get("/products/{id}", h.GetProduct)
}

func (h *Handlers) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.GetCategories(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *Handlers) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.GetProducts(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *Handlers) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.svc.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
