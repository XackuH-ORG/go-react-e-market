package catalog

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc/mock"
)

func TestCatalogService_GetCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := mock.NewMockQuerier(ctrl)
	svc := NewCatalogService(q)

	expectedCats := []repo.Category{
		{ID: uuid.New(), Name: "Electronics"},
		{ID: uuid.New(), Name: "Clothing"},
	}

	tests := []struct {
		name    string
		mock    func()
		want    []repo.Category
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				q.EXPECT().GetCategories(gomock.Any()).Return(expectedCats, nil)
			},
			want:    expectedCats,
			wantErr: false,
		},
		{
			name: "Error from DB",
			mock: func() {
				q.EXPECT().GetCategories(gomock.Any()).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.GetCategories(context.Background())
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCatalogService_GetProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := mock.NewMockQuerier(ctrl)
	svc := NewCatalogService(q)

	prodID1 := uuid.New()
	catID := uuid.New()
	skuID1 := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		want    []ProductResponse
		wantErr bool
	}{
		{
			name: "Success with skus",
			mock: func() {
				rows := []repo.GetProductsRow{
					{
						ProductID:   prodID1,
						CategoryID:  catID,
						ProductName: "Phone",
						Description: pgtype.Text{String: "A smartphone", Valid: true},
						SkuID:       pgtype.UUID{Bytes: skuID1, Valid: true},
						SkuCode:     pgtype.Text{String: "SKU-1", Valid: true},
						Price:       pgtype.Int4{Int32: 1000, Valid: true},
						Stock:       pgtype.Int4{Int32: 10, Valid: true},
						Attributes:  []byte(`{"color":"black"}`),
					},
					{
						ProductID:   prodID1,
						CategoryID:  catID,
						ProductName: "Phone",
						Description: pgtype.Text{String: "A smartphone", Valid: true},
						SkuID:       pgtype.UUID{Valid: false}, // this simulates a product with no more skus or left join nulls if handled that way
					},
				}
				q.EXPECT().GetProducts(gomock.Any()).Return(rows, nil)
			},
			want: []ProductResponse{
				{
					ID:          prodID1,
					CategoryID:  catID,
					Name:        "Phone",
					Description: "A smartphone",
					Skus: []SkuResponse{
						{
							ID:         skuID1,
							SkuCode:    "SKU-1",
							Price:      1000,
							Stock:      10,
							Attributes: []byte(`{"color":"black"}`),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				q.EXPECT().GetProducts(gomock.Any()).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.GetProducts(context.Background())
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Order of products in map iteration might be random, but here we only have 1
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCatalogService_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := mock.NewMockQuerier(ctrl)
	svc := NewCatalogService(q)

	prodID1 := uuid.New()
	catID := uuid.New()
	skuID1 := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		want    *ProductResponse
		wantErr bool
	}{
		{
			name: "Success",
			id:   prodID1,
			mock: func() {
				rows := []repo.GetProductRow{
					{
						ProductID:   prodID1,
						CategoryID:  catID,
						ProductName: "Phone",
						Description: pgtype.Text{String: "A smartphone", Valid: true},
						SkuID:       pgtype.UUID{Bytes: skuID1, Valid: true},
						SkuCode:     pgtype.Text{String: "SKU-1", Valid: true},
						Price:       pgtype.Int4{Int32: 1000, Valid: true},
						Stock:       pgtype.Int4{Int32: 10, Valid: true},
						Attributes:  []byte(`{"color":"black"}`),
					},
				}
				q.EXPECT().GetProduct(gomock.Any(), prodID1).Return(rows, nil)
			},
			want: &ProductResponse{
				ID:          prodID1,
				CategoryID:  catID,
				Name:        "Phone",
				Description: "A smartphone",
				Skus: []SkuResponse{
					{
						ID:         skuID1,
						SkuCode:    "SKU-1",
						Price:      1000,
						Stock:      10,
						Attributes: []byte(`{"color":"black"}`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			id:   prodID1,
			mock: func() {
				q.EXPECT().GetProduct(gomock.Any(), prodID1).Return([]repo.GetProductRow{}, nil)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "DB Error",
			id:   prodID1,
			mock: func() {
				q.EXPECT().GetProduct(gomock.Any(), prodID1).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.GetProduct(context.Background(), tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
