package products

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	product, err := h.service.CreateProduct(r.Context(), req.Name, req.Description, req.ImageURL)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(product)
}

type CreateSkuRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	SkuCode   string    `json:"sku_code"`
	Price     int32     `json:"price"`
	Stock     int32     `json:"stock"`
}

func (h *handler) CreateSku(w http.ResponseWriter, r *http.Request) {
	var req CreateSkuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Price <= 0 {
		http.Error(w, "Price must be greater than 0", http.StatusBadRequest)
		return
	}
	if req.Stock < 0 {
		http.Error(w, "Stock cannot be negative", http.StatusBadRequest)
		return
	}
	if req.SkuCode == "" {
		http.Error(w, "SkuCode is required", http.StatusBadRequest)
		return
	}

	sku, err := h.service.CreateSku(r.Context(), req.ProductID, req.SkuCode, req.Price, req.Stock)
	if err != nil {
		http.Error(w, "Failed to create SKU", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sku)
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	filters := SearchFilters{
		Limit:  10,
		Offset: 0,
	}

	query := r.URL.Query()
	filters.SearchQuery = query.Get("search")
	filters.InStock = query.Get("in_stock") == "true"

	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 32); err == nil && parsed > 0 {
			filters.Limit = int32(parsed)
		}
	}
	if o := query.Get("offset"); o != "" {
		if parsed, err := strconv.ParseInt(o, 10, 32); err == nil && parsed >= 0 {
			filters.Offset = int32(parsed)
		}
	}
	if p := query.Get("min_price"); p != "" {
		if parsed, err := strconv.ParseInt(p, 10, 32); err == nil {
			filters.MinPrice = int32(parsed)
		}
	}
	if p := query.Get("max_price"); p != "" {
		if parsed, err := strconv.ParseInt(p, 10, 32); err == nil {
			filters.MaxPrice = int32(parsed)
		}
	}

	products, err := h.service.ListProducts(r.Context(), filters)
	if err != nil {
		http.Error(w, "Failed to list products", http.StatusInternalServerError)
		return
	}

	if products == nil {
		products = make([]database.Product, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(products)
}

func (h *handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	productID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	details, err := h.service.GetProductDetails(r.Context(), productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(details)
}

type UpdateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	product, err := h.service.UpdateProduct(r.Context(), productID, req.Name, req.Description)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(product)
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduct(r.Context(), productID)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Ограничиваем размер до 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large or invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Путь для загрузки должен приходить из ENV, но для MVP берем локальную папку
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload dir", http.StatusInternalServerError)
		return
	}

	filename := productID.String() + filepath.Ext(header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	// Записываем URL (путь) в базу
	imageURL := "/uploads/" + filename
	product, err := h.service.UpdateProductImage(r.Context(), productID, imageURL)
	if err != nil {
		http.Error(w, "Failed to update DB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(product)
}
