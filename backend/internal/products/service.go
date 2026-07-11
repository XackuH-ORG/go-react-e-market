package products

import (
	"context"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductDetails struct {
	Product database.Product
	Skus    []database.Sku
}

// SearchFilters описывает параметры фильтрации
type SearchFilters struct {
	SearchQuery string
	MinPrice    int32
	MaxPrice    int32
	InStock     bool
	Limit       int32
	Offset      int32
}

type Service interface {
	CreateProduct(ctx context.Context, name, description, imageURL string) (database.Product, error)
	CreateSku(ctx context.Context, productID uuid.UUID, skuCode string, price int32, stock int32) (database.Sku, error)
	ListProducts(ctx context.Context, filters SearchFilters) ([]database.Product, error)
	GetProductDetails(ctx context.Context, id uuid.UUID) (ProductDetails, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, name, description string) (database.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	UpdateProductImage(ctx context.Context, id uuid.UUID, imageURL string) (database.Product, error)
}

// ProductStore defines the data access methods required by the products service
type ProductStore interface {
	CreateProduct(ctx context.Context, arg database.CreateProductParams) (database.Product, error)
	CreateSku(ctx context.Context, arg database.CreateSkuParams) (database.Sku, error)
	SearchProducts(ctx context.Context, arg database.SearchProductsParams) ([]database.Product, error)
	GetProductWithSkus(ctx context.Context, id uuid.UUID) ([]database.GetProductWithSkusRow, error)
	UpdateProduct(ctx context.Context, arg database.UpdateProductParams) (database.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	UpdateProductImage(ctx context.Context, arg database.UpdateProductImageParams) (database.Product, error)
}

type svc struct {
	db ProductStore
}

func NewService(db ProductStore) Service {
	return &svc{db: db}
}

// CreateProduct создает новый продукт
func (s *svc) CreateProduct(ctx context.Context, name, description, imageURL string) (database.Product, error) {
	desc := pgtype.Text{String: description, Valid: description != ""}
	img := pgtype.Text{String: imageURL, Valid: imageURL != ""}

	return s.db.CreateProduct(ctx, database.CreateProductParams{
		Name:        name,
		Description: desc,
		ImageUrl:    img,
	})
}

// CreateSku создает новый SKU
func (s *svc) CreateSku(ctx context.Context, productID uuid.UUID, skuCode string, price int32, stock int32) (database.Sku, error) {
	return s.db.CreateSku(ctx, database.CreateSkuParams{
		ProductID: productID,
		SkuCode:   skuCode,
		Price:     price,
		Stock:     stock,
	})
}

// ListProducts возвращает список продуктов с учетом фильтров
func (s *svc) ListProducts(ctx context.Context, filters SearchFilters) ([]database.Product, error) {
	return s.db.SearchProducts(ctx, database.SearchProductsParams{
		SearchQuery: filters.SearchQuery,
		MinPrice:    filters.MinPrice,
		MaxPrice:    filters.MaxPrice,
		InStock:     filters.InStock,
		LimitVal:    filters.Limit,
		OffsetVal:   filters.Offset,
	})
}

// GetProductByID is removed as we merged it

// GetSkusByProductID is removed as we merged it

// GetProductDetails получает продукт и его SKU в одном вызове (через JOIN)
func (s *svc) GetProductDetails(ctx context.Context, id uuid.UUID) (ProductDetails, error) {
	var details ProductDetails

	rows, err := s.db.GetProductWithSkus(ctx, id)
	if err != nil {
		return details, err
	}

	if len(rows) == 0 {
		return details, pgx.ErrNoRows // Имитируем отсутствие записи
	}

	// Заполняем информацию о продукте (она дублируется во всех строках)
	firstRow := rows[0]
	details.Product = database.Product{
		ID:          firstRow.PID,
		Name:        firstRow.PName,
		Description: firstRow.PDescription,
		ImageUrl:    firstRow.PImageUrl,
		CreatedAt:   firstRow.PCreatedAt,
		UpdatedAt:   firstRow.PUpdatedAt,
	}

	// Заполняем массив SKU
	for _, row := range rows {
		// Если s_id == nil (точнее Valid == false), значит у продукта еще нет SKU (результат LEFT JOIN)
		if !row.SID.Valid {
			continue
		}

		sku := database.Sku{
			ID:        row.SID.Bytes,
			ProductID: firstRow.PID,
			SkuCode:   row.SkuCode.String,
			Price:     row.Price.Int32,
			Stock:     row.Stock.Int32,
			CreatedAt: row.SCreatedAt.Time,
			UpdatedAt: row.SUpdatedAt.Time,
		}
		details.Skus = append(details.Skus, sku)
	}

	return details, nil
}

// UpdateProduct обновляет информацию о продукте
func (s *svc) UpdateProduct(ctx context.Context, id uuid.UUID, name, description string) (database.Product, error) {
	// sqlc ждет пустые строки или null-типы. Так как мы использовали COALESCE(NULLIF($2, '')), пустая строка игнорируется
	return s.db.UpdateProduct(ctx, database.UpdateProductParams{
		ID:      id,
		Column2: name,
		Column3: pgtype.Text{String: description, Valid: description != ""},
	})
}

// DeleteProduct мягко удаляет продукт
func (s *svc) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteProduct(ctx, id)
}

// UpdateProductImage привязывает URL картинки к товару
func (s *svc) UpdateProductImage(ctx context.Context, id uuid.UUID, imageURL string) (database.Product, error) {
	img := pgtype.Text{String: imageURL, Valid: imageURL != ""}
	return s.db.UpdateProductImage(ctx, database.UpdateProductImageParams{
		ID:       id,
		ImageUrl: img,
	})
}
