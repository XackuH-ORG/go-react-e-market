package orders_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc/mock"
	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/orders"
)

func dummyExecTx(m repo.Querier) func(ctx context.Context, fn func(repo.Querier) error) error {
	return func(ctx context.Context, fn func(repo.Querier) error) error {
		return fn(m)
	}
}

func TestOrdersService_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := orders.NewOrdersServiceForTest(m, dummyExecTx(m))

	ctx := context.Background()
	userID := uuid.New()
	expectedOrders := []repo.Order{{ID: uuid.New()}}

	t.Run("success", func(t *testing.T) {
		m.EXPECT().GetUserOrders(ctx, userID).Return(expectedOrders, nil)

		got, err := svc.GetUserOrders(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, expectedOrders, got)
	})

	t.Run("error", func(t *testing.T) {
		m.EXPECT().GetUserOrders(ctx, userID).Return(nil, errors.New("db error"))

		_, err := svc.GetUserOrders(ctx, userID)
		require.Error(t, err)
	})
}

func TestOrdersService_Checkout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := orders.NewOrdersServiceForTest(m, dummyExecTx(m))

	ctx := context.Background()
	userID := uuid.New()
	addr := "123 Main St"
	skuID := uuid.New()
	orderID := uuid.New()

	t.Run("success", func(t *testing.T) {
		items := []repo.GetCartItemsRow{
			{SkuID: skuID, Quantity: 2, Price: 100},
		}

		m.EXPECT().GetCartItems(ctx, userID).Return(items, nil)
		m.EXPECT().DecrementSkuStock(ctx, repo.DecrementSkuStockParams{
			Stock: 2,
			ID:    skuID,
		}).Return(repo.Sku{}, nil)

		m.EXPECT().CreateOrder(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, arg repo.CreateOrderParams) (repo.Order, error) {
			assert.Equal(t, userID, arg.UserID)
			assert.Equal(t, repo.DeliveryTypeDELIVERY, arg.DeliveryType)
			assert.Equal(t, addr, arg.DeliveryAddress.String)
			assert.Equal(t, int32(200), arg.TotalAmount)
			return repo.Order{ID: orderID}, nil
		})

		m.EXPECT().CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:         orderID,
			SkuID:           skuID,
			Quantity:        2,
			PriceAtPurchase: 100,
		}).Return(nil)

		m.EXPECT().ClearCart(ctx, userID).Return(nil)

		order, err := svc.Checkout(ctx, userID, repo.DeliveryTypeDELIVERY, &addr, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, orderID, order.ID)
	})

	t.Run("empty cart", func(t *testing.T) {
		m.EXPECT().GetCartItems(ctx, userID).Return([]repo.GetCartItemsRow{}, nil)

		_, err := svc.Checkout(ctx, userID, repo.DeliveryTypeDELIVERY, &addr, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cart is empty")
	})

	t.Run("not enough stock", func(t *testing.T) {
		items := []repo.GetCartItemsRow{
			{SkuID: skuID, Quantity: 2, Price: 100},
		}

		m.EXPECT().GetCartItems(ctx, userID).Return(items, nil)
		m.EXPECT().DecrementSkuStock(ctx, gomock.Any()).Return(repo.Sku{}, pgx.ErrNoRows)

		_, err := svc.Checkout(ctx, userID, repo.DeliveryTypeDELIVERY, &addr, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not enough stock")
	})

	t.Run("create order item error", func(t *testing.T) {
		items := []repo.GetCartItemsRow{
			{SkuID: skuID, Quantity: 2, Price: 100},
		}

		m.EXPECT().GetCartItems(ctx, userID).Return(items, nil)
		m.EXPECT().DecrementSkuStock(ctx, gomock.Any()).Return(repo.Sku{}, nil)

		m.EXPECT().CreateOrder(ctx, gomock.Any()).Return(repo.Order{ID: orderID}, nil)

		m.EXPECT().CreateOrderItem(ctx, gomock.Any()).Return(errors.New("db error"))

		_, err := svc.Checkout(ctx, userID, repo.DeliveryTypeDELIVERY, &addr, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create order item")
	})
}
