package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/admin"
	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc/mock"
)

func TestAdminService_GetAllOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	expectedOrders := []repo.Order{{ID: uuid.New()}, {ID: uuid.New()}}

	tests := []struct {
		name    string
		mock    func()
		want    []repo.Order
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				m.EXPECT().GetAllOrders(ctx).Return(expectedOrders, nil)
			},
			want:    expectedOrders,
			wantErr: false,
		},
		{
			name: "error",
			mock: func() {
				m.EXPECT().GetAllOrders(ctx).Return(nil, pgx.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.GetAllOrders(ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAdminService_SearchOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	expectedOrders := []repo.Order{{ID: uuid.New()}}

	tests := []struct {
		name    string
		index   string
		mock    func(index string)
		want    []repo.Order
		wantErr bool
	}{
		{
			name:  "success",
			index: "1234",
			mock: func(index string) {
				m.EXPECT().SearchOrders(ctx, pgtype.Text{String: index, Valid: true}).Return(expectedOrders, nil)
			},
			want:    expectedOrders,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.index)
			got, err := svc.SearchOrders(ctx, tt.index)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAdminService_UpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	orderID := uuid.New()

	tests := []struct {
		name    string
		status  repo.OrderStatus
		mock    func(status repo.OrderStatus)
		wantErr bool
	}{
		{
			name:   "success",
			status: repo.OrderStatusDELIVERED,
			mock: func(status repo.OrderStatus) {
				m.EXPECT().UpdateOrderStatus(ctx, repo.UpdateOrderStatusParams{
					Status: status,
					ID:     orderID,
				}).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.status)
			err := svc.UpdateOrderStatus(ctx, orderID, tt.status)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAdminService_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	expectedUsers := []repo.User{{ID: uuid.New()}}

	tests := []struct {
		name    string
		mock    func()
		want    []repo.User
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				m.EXPECT().GetUsers(ctx).Return(expectedUsers, nil)
			},
			want:    expectedUsers,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.GetUsers(ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAdminService_UpdateUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name    string
		role    repo.UserRole
		mock    func(role repo.UserRole)
		wantErr bool
	}{
		{
			name: "success",
			role: repo.UserRoleADMIN,
			mock: func(role repo.UserRole) {
				m.EXPECT().UpdateUserRole(ctx, repo.UpdateUserRoleParams{
					Role: role,
					ID:   userID,
				}).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.role)
			err := svc.UpdateUserRole(ctx, userID, tt.role)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAdminService_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	categoryID := uuid.New()
	name := "Product"
	desc := "Desc"
	expectedProduct := repo.Product{ID: uuid.New(), Name: name}

	tests := []struct {
		name    string
		desc    *string
		mock    func(desc *string)
		want    repo.Product
		wantErr bool
	}{
		{
			name: "success with desc",
			desc: &desc,
			mock: func(d *string) {
				m.EXPECT().CreateProduct(ctx, repo.CreateProductParams{
					CategoryID:  categoryID,
					Name:        name,
					Description: pgtype.Text{String: *d, Valid: true},
				}).Return(expectedProduct, nil)
			},
			want:    expectedProduct,
			wantErr: false,
		},
		{
			name: "success without desc",
			desc: nil,
			mock: func(d *string) {
				m.EXPECT().CreateProduct(ctx, repo.CreateProductParams{
					CategoryID:  categoryID,
					Name:        name,
					Description: pgtype.Text{},
				}).Return(expectedProduct, nil)
			},
			want:    expectedProduct,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.desc)
			got, err := svc.CreateProduct(ctx, categoryID, name, tt.desc)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAdminService_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	productID := uuid.New()
	categoryID := uuid.New()
	name := "Product"
	desc := "Desc"

	tests := []struct {
		name    string
		desc    *string
		mock    func(desc *string)
		wantErr bool
	}{
		{
			name: "success",
			desc: &desc,
			mock: func(d *string) {
				m.EXPECT().UpdateProduct(ctx, repo.UpdateProductParams{
					CategoryID:  categoryID,
					Name:        name,
					Description: pgtype.Text{String: *d, Valid: true},
					ID:          productID,
				}).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.desc)
			err := svc.UpdateProduct(ctx, productID, categoryID, name, tt.desc)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAdminService_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	productID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				m.EXPECT().SoftDeleteProduct(ctx, productID).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := svc.DeleteProduct(ctx, productID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAdminService_CreateSku(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	productID := uuid.New()
	skuCode := "SKU-1"
	attr := []byte(`{"color":"red"}`)
	price := int32(100)
	stock := int32(10)
	expectedSku := repo.Sku{ID: uuid.New(), SkuCode: skuCode}

	tests := []struct {
		name    string
		mock    func()
		want    repo.Sku
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				m.EXPECT().CreateSku(ctx, repo.CreateSkuParams{
					ProductID:  productID,
					SkuCode:    skuCode,
					Attributes: attr,
					Price:      price,
					Stock:      stock,
				}).Return(expectedSku, nil)
			},
			want:    expectedSku,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.CreateSku(ctx, productID, skuCode, attr, price, stock)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAdminService_UpdateSku(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	skuID := uuid.New()
	skuCode := "SKU-1"
	attr := []byte(`{"color":"red"}`)
	price := int32(100)
	stock := int32(10)

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				m.EXPECT().UpdateSku(ctx, repo.UpdateSkuParams{
					ID:         skuID,
					SkuCode:    skuCode,
					Attributes: attr,
					Price:      price,
					Stock:      stock,
				}).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := svc.UpdateSku(ctx, skuID, skuCode, attr, price, stock)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAdminService_DeleteSku(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockQuerier(ctrl)
	svc := admin.NewAdminService(nil, m)

	ctx := context.Background()
	skuID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				m.EXPECT().SoftDeleteSku(ctx, skuID).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := svc.DeleteSku(ctx, skuID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
