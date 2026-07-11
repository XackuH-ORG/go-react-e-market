package cart

import (
	"context"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/google/uuid"
)

// Store defines the interface for database operations related to the cart.
type Store interface {
	AddToCart(ctx context.Context, arg database.AddToCartParams) (database.CartItem, error)
	GetCartItems(ctx context.Context, userID uuid.UUID) ([]database.GetCartItemsRow, error)
	RemoveFromCart(ctx context.Context, arg database.RemoveFromCartParams) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

type Service interface {
	AddToCart(ctx context.Context, userID, skuID uuid.UUID, quantity int32) (database.CartItem, error)
	GetCart(ctx context.Context, userID uuid.UUID) ([]database.GetCartItemsRow, error)
	RemoveFromCart(ctx context.Context, userID, skuID uuid.UUID) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{
		store: store,
	}
}

func (s *service) AddToCart(ctx context.Context, userID, skuID uuid.UUID, quantity int32) (database.CartItem, error) {
	return s.store.AddToCart(ctx, database.AddToCartParams{
		UserID:   userID,
		SkuID:    skuID,
		Quantity: quantity,
	})
}

func (s *service) GetCart(ctx context.Context, userID uuid.UUID) ([]database.GetCartItemsRow, error) {
	return s.store.GetCartItems(ctx, userID)
}

func (s *service) RemoveFromCart(ctx context.Context, userID, skuID uuid.UUID) error {
	return s.store.RemoveFromCart(ctx, database.RemoveFromCartParams{
		UserID: userID,
		SkuID:  skuID,
	})
}

func (s *service) ClearCart(ctx context.Context, userID uuid.UUID) error {
	return s.store.ClearCart(ctx, userID)
}
