package cart

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type CartService interface {
	AddToCart(ctx context.Context, userID, skuID uuid.UUID, quantity int32) error
	GetCart(ctx context.Context, userID uuid.UUID) ([]repo.GetCartItemsRow, error)
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

type cartService struct {
	q repo.Querier
}

func NewCartService(q repo.Querier) CartService {
	return &cartService{q: q}
}

func (s *cartService) AddToCart(ctx context.Context, userID, skuID uuid.UUID, quantity int32) error {
	return s.q.AddToCart(ctx, repo.AddToCartParams{
		UserID:   userID,
		SkuID:    skuID,
		Quantity: quantity,
	})
}

func (s *cartService) GetCart(ctx context.Context, userID uuid.UUID) ([]repo.GetCartItemsRow, error) {
	return s.q.GetCartItems(ctx, userID)
}

func (s *cartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	return s.q.ClearCart(ctx, userID)
}
