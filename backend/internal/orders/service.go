package orders

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type OrdersService interface {
	Checkout(ctx context.Context, userID uuid.UUID, deliveryType repo.DeliveryType, deliveryAddress *string, pickupPointID *uuid.UUID, promoCodeID *uuid.UUID) (repo.Order, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID) ([]repo.Order, error)
}

type ordersService struct {
	db     *pgxpool.Pool
	q      repo.Querier
	execTx func(ctx context.Context, fn func(repo.Querier) error) error
}

func defaultExecTx(db *pgxpool.Pool) func(ctx context.Context, fn func(repo.Querier) error) error {
	return func(ctx context.Context, fn func(repo.Querier) error) error {
		tx, err := db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		qtx := repo.New(db).WithTx(tx)
		err = fn(qtx)
		if err != nil {
			return err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		return nil
	}
}

func NewOrdersService(db *pgxpool.Pool) OrdersService {
	return &ordersService{
		db:     db,
		q:      repo.New(db),
		execTx: defaultExecTx(db),
	}
}

func NewOrdersServiceForTest(q repo.Querier, execTx func(ctx context.Context, fn func(repo.Querier) error) error) OrdersService {
	return &ordersService{
		q:      q,
		execTx: execTx,
	}
}

func generatePublicNumber() string {
	// Simple generator for public number
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("ORD-%06d-%04d", r.Intn(1000000), r.Intn(10000))
}

func (s *ordersService) Checkout(ctx context.Context, userID uuid.UUID, deliveryType repo.DeliveryType, deliveryAddress *string, pickupPointID *uuid.UUID, promoCodeID *uuid.UUID) (repo.Order, error) {
	var createdOrder repo.Order

	err := s.execTx(ctx, func(qtx repo.Querier) error {
		// Get cart items
		items, err := qtx.GetCartItems(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get cart items: %w", err)
		}

		if len(items) == 0 {
			return errors.New("cart is empty")
		}

		var totalAmount int32 = 0

		for _, item := range items {
			// Decrement stock
			_, err := qtx.DecrementSkuStock(ctx, repo.DecrementSkuStockParams{
				Stock: item.Quantity,
				ID:    item.SkuID,
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return fmt.Errorf("not enough stock for sku %s", item.SkuID)
				}
				return fmt.Errorf("failed to decrement stock: %w", err)
			}

			totalAmount += item.Price * item.Quantity
		}

		// Create order
		// Adjust address and pickup point
		var addr pgtype.Text
		if deliveryAddress != nil {
			addr = pgtype.Text{String: *deliveryAddress, Valid: true}
		}
		var pp pgtype.UUID
		if pickupPointID != nil {
			pp = pgtype.UUID{Bytes: *pickupPointID, Valid: true}
		}
		var promo pgtype.UUID
		if promoCodeID != nil {
			promo = pgtype.UUID{Bytes: *promoCodeID, Valid: true}
		}

		createdOrder, err = qtx.CreateOrder(ctx, repo.CreateOrderParams{
			PublicNumber:    generatePublicNumber(),
			UserID:          userID,
			DeliveryType:    deliveryType,
			DeliveryAddress: addr,
			PickupPointID:   pp,
			TotalAmount:     totalAmount,
			PromoCodeID:     promo,
		})
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Create order items
		for _, item := range items {
			err := qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
				OrderID:         createdOrder.ID,
				SkuID:           item.SkuID,
				Quantity:        item.Quantity,
				PriceAtPurchase: item.Price,
			})
			if err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}

		// Clear cart
		err = qtx.ClearCart(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to clear cart: %w", err)
		}

		return nil
	})

	return createdOrder, err
}

func (s *ordersService) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]repo.Order, error) {
	return s.q.GetUserOrders(ctx, userID)
}
