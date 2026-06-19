package auth_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/auth"
	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	mock_repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc/mock"
)

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_repo.NewMockQuerier(ctrl)
	jwtSecret := []byte("secret")
	service := auth.NewAuthService(mockQuerier, jwtSecret)

	email := "test@example.com"
	password := "password123"

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
		checkResult   func(t *testing.T, user *repo.User, err error)
	}{
		{
			name: "Success",
			setupMock: func() {
				mockQuerier.EXPECT().GetUserByEmail(gomock.Any(), email).Return(repo.User{}, sql.ErrNoRows)

				mockQuerier.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) {
					assert.Equal(t, email, arg.Email)
					assert.Equal(t, repo.UserRoleCUSTOMER, arg.Role)
					// Let's not assert role strictly here, we will just return it.
					err := bcrypt.CompareHashAndPassword([]byte(arg.PasswordHash), []byte(password))
					assert.NoError(t, err)

					return repo.User{
						ID:           uuid.New(),
						Email:        arg.Email,
						PasswordHash: arg.PasswordHash,
						Role:         arg.Role,
					}, nil
				})
			},
			checkResult: func(t *testing.T, user *repo.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, email, user.Email)
			},
		},
		{
			name: "User Already Exists",
			setupMock: func() {
				mockQuerier.EXPECT().GetUserByEmail(gomock.Any(), email).Return(repo.User{
					ID:    uuid.New(),
					Email: email,
				}, nil)
			},
			checkResult: func(t *testing.T, user *repo.User, err error) {
				require.ErrorIs(t, err, auth.ErrUserAlreadyExists)
				require.Nil(t, user)
			},
		},
		{
			name: "Database Error on Create",
			setupMock: func() {
				mockQuerier.EXPECT().GetUserByEmail(gomock.Any(), email).Return(repo.User{}, sql.ErrNoRows)
				mockQuerier.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(repo.User{}, sql.ErrConnDone)
			},
			checkResult: func(t *testing.T, user *repo.User, err error) {
				require.ErrorIs(t, err, sql.ErrConnDone)
				require.Nil(t, user)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			user, err := service.Register(context.Background(), email, password)
			tc.checkResult(t, user, err)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_repo.NewMockQuerier(ctrl)
	jwtSecret := []byte("secret")
	service := auth.NewAuthService(mockQuerier, jwtSecret)

	email := "test@example.com"
	password := "password123"

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
		checkResult   func(t *testing.T, user *repo.User, token string, err error)
	}{
		{
			name: "Success",
			setupMock: func() {
				hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				mockQuerier.EXPECT().GetUserByEmail(gomock.Any(), email).Return(repo.User{
					ID:           uuid.New(),
					Email:        email,
					PasswordHash: string(hash),
					Role:         "User", // we just provide whatever string
				}, nil)
			},
			checkResult: func(t *testing.T, user *repo.User, token string, err error) {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.NotEmpty(t, token)
			},
		},
		{
			name: "User Not Found",
			setupMock: func() {
				mockQuerier.EXPECT().GetUserByEmail(gomock.Any(), email).Return(repo.User{}, sql.ErrNoRows)
			},
			checkResult: func(t *testing.T, user *repo.User, token string, err error) {
				require.ErrorIs(t, err, auth.ErrInvalidCredentials)
				require.Empty(t, token)
				require.Nil(t, user)
			},
		},
		{
			name: "Invalid Password",
			setupMock: func() {
				hash, _ := bcrypt.GenerateFromPassword([]byte("wrongpassword"), bcrypt.DefaultCost)
				mockQuerier.EXPECT().GetUserByEmail(gomock.Any(), email).Return(repo.User{
					ID:           uuid.New(),
					Email:        email,
					PasswordHash: string(hash),
					Role:         "User",
				}, nil)
			},
			checkResult: func(t *testing.T, user *repo.User, token string, err error) {
				require.ErrorIs(t, err, auth.ErrInvalidCredentials)
				require.Empty(t, token)
				require.Nil(t, user)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			user, token, err := service.Login(context.Background(), email, password)
			tc.checkResult(t, user, token, err)
		})
	}
}
