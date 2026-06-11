package catalog

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type SkuResponse struct {
	ID         uuid.UUID `json:"id"`
	SkuCode    string    `json:"sku_code"`
	Price      int32     `json:"price"`
	Stock      int32     `json:"stock"`
	Attributes []byte    `json:"attributes"`
}

type ProductResponse struct {
	ID          uuid.UUID     `json:"id"`
	CategoryID  uuid.UUID     `json:"category_id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Skus        []SkuResponse `json:"skus"`
}

type CatalogService interface {
	GetCategories(ctx context.Context) ([]repo.Category, error)
	GetProducts(ctx context.Context) ([]ProductResponse, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*ProductResponse, error)
}

type catalogService struct {
	q repo.Querier
}

func NewCatalogService(q repo.Querier) CatalogService {
	return &catalogService{q: q}
}

func (s *catalogService) GetCategories(ctx context.Context) ([]repo.Category, error) {
	return s.q.GetCategories(ctx)
}

func (s *catalogService) GetProducts(ctx context.Context) ([]ProductResponse, error) {
	rows, err := s.q.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	productMap := make(map[uuid.UUID]*ProductResponse)
	for _, row := range rows {
		if _, ok := productMap[row.ProductID]; !ok {
			productMap[row.ProductID] = &ProductResponse{
				ID:          row.ProductID,
				CategoryID:  row.CategoryID,
				Name:        row.ProductName,
				Description: row.Description.String,
				Skus:        []SkuResponse{},
			}
		}

		if row.SkuID.Valid {
			sku := SkuResponse{
				ID:         row.SkuID.Bytes,
				SkuCode:    row.SkuCode.String,
				Price:      row.Price.Int32,
				Stock:      row.Stock.Int32,
				Attributes: row.Attributes,
			}
			productMap[row.ProductID].Skus = append(productMap[row.ProductID].Skus, sku)
		}
	}

	products := make([]ProductResponse, 0, len(productMap))
	for _, p := range productMap {
		products = append(products, *p)
	}

	return products, nil
}

func (s *catalogService) GetProduct(ctx context.Context, id uuid.UUID) (*ProductResponse, error) {
	rows, err := s.q.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil // Not found
	}

	product := &ProductResponse{
		ID:          rows[0].ProductID,
		CategoryID:  rows[0].CategoryID,
		Name:        rows[0].ProductName,
		Description: rows[0].Description.String,
		Skus:        []SkuResponse{},
	}

	for _, row := range rows {
		if row.SkuID.Valid {
			sku := SkuResponse{
				ID:         row.SkuID.Bytes,
				SkuCode:    row.SkuCode.String,
				Price:      row.Price.Int32,
				Stock:      row.Stock.Int32,
				Attributes: row.Attributes,
			}
			product.Skus = append(product.Skus, sku)
		}
	}

	return product, nil
}
