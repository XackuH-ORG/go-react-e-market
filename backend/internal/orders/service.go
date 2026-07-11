package orders

import (
	"context"
	"fmt"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateOrderRequest struct {
	DeliveryType    string    `json:"delivery_type"`
	DeliveryAddress string    `json:"delivery_address,omitempty"`
	PickupPointID   uuid.UUID `json:"pickup_point_id,omitempty"`
}

type Service interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (database.Order, error)
	GetOrders(ctx context.Context, userID uuid.UUID) ([]database.Order, error)
	GetAdminOrders(ctx context.Context, searchIndex string) ([]database.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) (database.Order, error)
}

type service struct {
	pool *pgxpool.Pool
	db   *database.Queries
}

func NewService(pool *pgxpool.Pool) Service {
	return &service{
		pool: pool,
		db:   database.New(pool),
	}
}

func (s *service) CreateOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (database.Order, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return database.Order{}, err
	}
	defer tx.Rollback(ctx)

	q := s.db.WithTx(tx)

	items, err := q.GetCartItems(ctx, userID)
	if err != nil {
		return database.Order{}, err
	}
	if len(items) == 0 {
		return database.Order{}, fmt.Errorf("cart is empty")
	}

	var totalAmount int32 = 0

	for _, item := range items {
		_, err := q.DecrementSkuStock(ctx, database.DecrementSkuStockParams{
			Stock: item.Quantity,
			ID:    item.SkuID,
		})
		if err != nil {
			return database.Order{}, fmt.Errorf("failed to reserve stock for sku %v: %w", item.SkuID, err)
		}
		totalAmount += item.Price * item.Quantity
	}

	publicNumber := fmt.Sprintf("ORD-%s", uuid.NewString()[:8])

	var deliveryAddress pgtype.Text
	if req.DeliveryAddress != "" {
		deliveryAddress = pgtype.Text{String: req.DeliveryAddress, Valid: true}
	}

	var pickupPointID pgtype.UUID
	if req.PickupPointID != uuid.Nil {
		pickupPointID = pgtype.UUID{Bytes: req.PickupPointID, Valid: true}
	}

	order, err := q.CreateOrder(ctx, database.CreateOrderParams{
		PublicNumber:    publicNumber,
		UserID:          userID,
		DeliveryType:    database.DeliveryType(req.DeliveryType),
		DeliveryAddress: deliveryAddress,
		PickupPointID:   pickupPointID,
		TotalAmount:     totalAmount,
		PromoCodeID:     pgtype.UUID{},
	})
	if err != nil {
		return database.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range items {
		err := q.CreateOrderItem(ctx, database.CreateOrderItemParams{
			OrderID:         order.ID,
			SkuID:           item.SkuID,
			Quantity:        item.Quantity,
			PriceAtPurchase: item.Price,
		})
		if err != nil {
			return database.Order{}, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	err = q.ClearCart(ctx, userID)
	if err != nil {
		return database.Order{}, fmt.Errorf("failed to clear cart: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return database.Order{}, err
	}

	return order, nil
}

func (s *service) GetOrders(ctx context.Context, userID uuid.UUID) ([]database.Order, error) {
	return s.db.GetUserOrders(ctx, userID)
}

func (s *service) GetAdminOrders(ctx context.Context, searchIndex string) ([]database.Order, error) {
	if searchIndex != "" {
		return s.db.SearchOrders(ctx, pgtype.Text{String: searchIndex, Valid: true})
	}
	// Fallback to a limited list of orders for admin if search is empty
	return []database.Order{}, nil
}

func (s *service) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) (database.Order, error) {
	return s.db.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID:     orderID,
		Status: database.OrderStatus(status),
	})
}
