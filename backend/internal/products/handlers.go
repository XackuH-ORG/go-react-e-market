package products

import (
	"log/slog"
	"net/http"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	err := h.service.ListProducts(r.Context())
	if err != nil {
		slog.Error("Ошибка при получении списка продуктов", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := struct {
		Products []string `json:"products"`
	}{}

	json.Write(w, http.StatusOK, products)

}
