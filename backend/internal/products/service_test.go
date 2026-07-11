package products

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/products/mocks"
)

func TestService_CreateProduct(t *testing.T) {
	t.Parallel()

	expectedProduct := database.Product{
		ID:          uuid.New(),
		Name:        "Test Product",
		Description: pgtype.Text{String: "Desc", Valid: true},
		ImageUrl:    pgtype.Text{String: "http://img.com", Valid: true},
	}

	tests := []struct {
		name        string
		productName string
		description string
		imageURL    string
		mockSetup   func(m *mocks.ProductStore)
		wantErr     error
		wantResult  database.Product
	}{
		{
			name:        "success",
			productName: "Test Product",
			description: "Desc",
			imageURL:    "http://img.com",
			mockSetup: func(m *mocks.ProductStore) {
				m.EXPECT().
					CreateProduct(mock.Anything, mock.Anything).
					Return(expectedProduct, nil).
					Once()
			},
			wantErr:    nil,
			wantResult: expectedProduct,
		},
		{
			name:        "db error",
			productName: "Test Product",
			description: "Desc",
			imageURL:    "http://img.com",
			mockSetup: func(m *mocks.ProductStore) {
				m.EXPECT().
					CreateProduct(mock.Anything, mock.Anything).
					Return(database.Product{}, errors.New("db error")).
					Once()
			},
			wantErr:    errors.New("db error"),
			wantResult: database.Product{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockStore := mocks.NewProductStore(t)
			tt.mockSetup(mockStore)
			svc := NewService(mockStore)

			result, err := svc.CreateProduct(context.Background(), tt.productName, tt.description, tt.imageURL)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, result)
			}
		})
	}
}

func TestService_GetProductDetails(t *testing.T) {
	t.Parallel()
	productID := uuid.New()
	expectedSkus := []database.Sku{{ID: uuid.New(), ProductID: productID, SkuCode: "SKU-1"}}

	tests := []struct {
		name       string
		productID  uuid.UUID
		mockSetup  func(m *mocks.ProductStore)
		wantErr    error
		wantResult ProductDetails
	}{
		{
			name:      "success",
			productID: productID,
			mockSetup: func(m *mocks.ProductStore) {
				rows := []database.GetProductWithSkusRow{
					{
						PID:     productID,
						PName:   "Test",
						SID:     pgtype.UUID{Bytes: expectedSkus[0].ID, Valid: true},
						SkuCode: pgtype.Text{String: "SKU-1", Valid: true},
					},
				}
				m.EXPECT().GetProductWithSkus(mock.Anything, productID).Return(rows, nil).Once()
			},
			wantErr: nil,
			wantResult: ProductDetails{
				Product: database.Product{ID: productID, Name: "Test"},
				Skus: []database.Sku{
					{ID: expectedSkus[0].ID, ProductID: productID, SkuCode: "SKU-1"},
				},
			},
		},
		{
			name:      "product not found",
			productID: productID,
			mockSetup: func(m *mocks.ProductStore) {
				m.EXPECT().GetProductWithSkus(mock.Anything, productID).Return(nil, nil).Once()
			},
			wantErr:    pgx.ErrNoRows,
			wantResult: ProductDetails{},
		},
		{
			name:      "db error",
			productID: productID,
			mockSetup: func(m *mocks.ProductStore) {
				m.EXPECT().GetProductWithSkus(mock.Anything, productID).Return(nil, errors.New("db error")).Once()
			},
			wantErr:    errors.New("db error"),
			wantResult: ProductDetails{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockStore := mocks.NewProductStore(t)
			tt.mockSetup(mockStore)
			svc := NewService(mockStore)

			details, err := svc.GetProductDetails(context.Background(), tt.productID)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, details)
			}
		})
	}
}

func TestService_ListProducts(t *testing.T) {
	t.Parallel()

	expectedProducts := []database.Product{{ID: uuid.New(), Name: "Apple iPad"}}
	filters := SearchFilters{SearchQuery: "apple", Limit: 10}

	tests := []struct {
		name       string
		filters    SearchFilters
		mockSetup  func(m *mocks.ProductStore)
		wantErr    error
		wantResult []database.Product
	}{
		{
			name:    "success",
			filters: filters,
			mockSetup: func(m *mocks.ProductStore) {
				m.EXPECT().SearchProducts(mock.Anything, mock.Anything).Return(expectedProducts, nil).Once()
			},
			wantErr:    nil,
			wantResult: expectedProducts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockStore := mocks.NewProductStore(t)
			tt.mockSetup(mockStore)
			svc := NewService(mockStore)

			result, err := svc.ListProducts(context.Background(), tt.filters)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, result)
			}
		})
	}
}

func TestService_UpdateProduct(t *testing.T) {
	t.Parallel()
	productID := uuid.New()
	expectedProduct := database.Product{ID: productID, Name: "New Name"}

	tests := []struct {
		name        string
		productID   uuid.UUID
		productName string
		description string
		mockSetup   func(m *mocks.ProductStore)
		wantErr     error
		wantResult  database.Product
	}{
		{
			name:        "success",
			productID:   productID,
			productName: "New Name",
			description: "New Desc",
			mockSetup: func(m *mocks.ProductStore) {
				m.EXPECT().UpdateProduct(mock.Anything, mock.Anything).Return(expectedProduct, nil).Once()
			},
			wantErr:    nil,
			wantResult: expectedProduct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockStore := mocks.NewProductStore(t)
			tt.mockSetup(mockStore)
			svc := NewService(mockStore)

			result, err := svc.UpdateProduct(context.Background(), tt.productID, tt.productName, tt.description)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, result)
			}
		})
	}
}

func TestService_DeleteProduct(t *testing.T) {
	t.Parallel()
	productID := uuid.New()

	tests := []struct {
		name      string
		productID uuid.UUID
		mockSetup func(m *mocks.ProductStore)
		wantErr   error
	}{
		{
			name:      "success",
			productID: productID,
			mockSetup: func(m *mocks.ProductStore) {
				m.EXPECT().DeleteProduct(mock.Anything, productID).Return(nil).Once()
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockStore := mocks.NewProductStore(t)
			tt.mockSetup(mockStore)
			svc := NewService(mockStore)

			err := svc.DeleteProduct(context.Background(), tt.productID)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
