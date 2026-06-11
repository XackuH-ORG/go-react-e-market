package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc/mock"
)

func TestCartService_AddToCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := mock.NewMockQuerier(ctrl)
	svc := NewCartService(q)

	userID := uuid.New()
	skuID := uuid.New()

	tests := []struct {
		name     string
		quantity int32
		mock     func()
		wantErr  bool
	}{
		{
			name:     "Success",
			quantity: 2,
			mock: func() {
				q.EXPECT().AddToCart(gomock.Any(), repo.AddToCartParams{
					UserID:   userID,
					SkuID:    skuID,
					Quantity: 2,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "DB Error",
			quantity: 1,
			mock: func() {
				q.EXPECT().AddToCart(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := svc.AddToCart(context.Background(), userID, skuID, tt.quantity)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCartService_GetCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := mock.NewMockQuerier(ctrl)
	svc := NewCartService(q)

	userID := uuid.New()
	expectedItems := []repo.GetCartItemsRow{
		{SkuID: uuid.New(), Quantity: 1},
		{SkuID: uuid.New(), Quantity: 3},
	}

	tests := []struct {
		name    string
		mock    func()
		want    []repo.GetCartItemsRow
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				q.EXPECT().GetCartItems(gomock.Any(), userID).Return(expectedItems, nil)
			},
			want:    expectedItems,
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				q.EXPECT().GetCartItems(gomock.Any(), userID).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.GetCart(context.Background(), userID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCartService_ClearCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := mock.NewMockQuerier(ctrl)
	svc := NewCartService(q)

	userID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				q.EXPECT().ClearCart(gomock.Any(), userID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				q.EXPECT().ClearCart(gomock.Any(), userID).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := svc.ClearCart(context.Background(), userID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
