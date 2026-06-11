package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type AdminService interface {
	GetAllOrders(ctx context.Context) ([]repo.Order, error)
	SearchOrders(ctx context.Context, index string) ([]repo.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status repo.OrderStatus) error
	
	GetUsers(ctx context.Context) ([]repo.User, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role repo.UserRole) error

	CreateProduct(ctx context.Context, categoryID uuid.UUID, name string, description *string) (repo.Product, error)
	UpdateProduct(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID, name string, description *string) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error

	CreateSku(ctx context.Context, productID uuid.UUID, skuCode string, attributes []byte, price int32, stock int32) (repo.Sku, error)
	UpdateSku(ctx context.Context, skuID uuid.UUID, skuCode string, attributes []byte, price int32, stock int32) error
	DeleteSku(ctx context.Context, skuID uuid.UUID) error
}

type adminService struct {
	db *pgxpool.Pool
	q  repo.Querier
}

func NewAdminService(db *pgxpool.Pool, q repo.Querier) AdminService {
	return &adminService{db: db, q: q}
}

func (s *adminService) GetAllOrders(ctx context.Context) ([]repo.Order, error) {
	return s.q.GetAllOrders(ctx)
}

func (s *adminService) SearchOrders(ctx context.Context, index string) ([]repo.Order, error) {
	var searchIndex pgtype.Text
	searchIndex.String = index
	searchIndex.Valid = true
	return s.q.SearchOrders(ctx, searchIndex)
}

func (s *adminService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status repo.OrderStatus) error {
	return s.q.UpdateOrderStatus(ctx, repo.UpdateOrderStatusParams{
		Status: status,
		ID:     orderID,
	})
}

func (s *adminService) GetUsers(ctx context.Context) ([]repo.User, error) {
	return s.q.GetUsers(ctx)
}

func (s *adminService) UpdateUserRole(ctx context.Context, userID uuid.UUID, role repo.UserRole) error {
	return s.q.UpdateUserRole(ctx, repo.UpdateUserRoleParams{
		Role: role,
		ID:   userID,
	})
}

func (s *adminService) CreateProduct(ctx context.Context, categoryID uuid.UUID, name string, description *string) (repo.Product, error) {
	var desc pgtype.Text
	if description != nil {
		desc = pgtype.Text{String: *description, Valid: true}
	}
	return s.q.CreateProduct(ctx, repo.CreateProductParams{
		CategoryID:  categoryID,
		Name:        name,
		Description: desc,
	})
}

func (s *adminService) UpdateProduct(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID, name string, description *string) error {
	var desc pgtype.Text
	if description != nil {
		desc = pgtype.Text{String: *description, Valid: true}
	}
	return s.q.UpdateProduct(ctx, repo.UpdateProductParams{
		CategoryID:  categoryID,
		Name:        name,
		Description: desc,
		ID:          productID,
	})
}

func (s *adminService) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return s.q.SoftDeleteProduct(ctx, productID)
}

func (s *adminService) CreateSku(ctx context.Context, productID uuid.UUID, skuCode string, attributes []byte, price int32, stock int32) (repo.Sku, error) {
	return s.q.CreateSku(ctx, repo.CreateSkuParams{
		ProductID:  productID,
		SkuCode:    skuCode,
		Attributes: attributes,
		Price:      price,
		Stock:      stock,
	})
}

func (s *adminService) UpdateSku(ctx context.Context, skuID uuid.UUID, skuCode string, attributes []byte, price int32, stock int32) error {
	return s.q.UpdateSku(ctx, repo.UpdateSkuParams{
		SkuCode:    skuCode,
		Attributes: attributes,
		Price:      price,
		Stock:      stock,
		ID:         skuID,
	})
}

func (s *adminService) DeleteSku(ctx context.Context, skuID uuid.UUID) error {
	return s.q.SoftDeleteSku(ctx, skuID)
}
