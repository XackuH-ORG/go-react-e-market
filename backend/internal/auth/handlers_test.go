package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/auth"
	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type mockAuthService struct {
	registerFunc func(ctx context.Context, email, password string) (*repo.User, error)
	loginFunc    func(ctx context.Context, email, password string) (string, error)
}

func (m *mockAuthService) Register(ctx context.Context, email, password string) (*repo.User, error) {
	return m.registerFunc(ctx, email, password)
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	return m.loginFunc(ctx, email, password)
}

func TestHandlers_Register(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		mockService    *mockAuthService
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Success",
			body: auth.RegisterRequest{Email: "test@example.com", Password: "pass"},
			mockService: &mockAuthService{
				registerFunc: func(ctx context.Context, email, password string) (*repo.User, error) {
					return &repo.User{
						ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Email: email,
						Role:  repo.UserRoleCUSTOMER,
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":    "00000000-0000-0000-0000-000000000001",
				"email": "test@example.com",
				"role":  string(repo.UserRoleCUSTOMER),
			},
		},
		{
			name: "Invalid Body",
			body: "invalid-json",
			mockService: &mockAuthService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing Email",
			body: auth.RegisterRequest{Password: "pass"},
			mockService: &mockAuthService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "User Already Exists",
			body: auth.RegisterRequest{Email: "test@example.com", Password: "pass"},
			mockService: &mockAuthService{
				registerFunc: func(ctx context.Context, email, password string) (*repo.User, error) {
					return nil, auth.ErrUserAlreadyExists
				},
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "Internal Error",
			body: auth.RegisterRequest{Email: "test@example.com", Password: "pass"},
			mockService: &mockAuthService{
				registerFunc: func(ctx context.Context, email, password string) (*repo.User, error) {
					return nil, errors.New("db error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := auth.NewHandlers(tc.mockService)

			var buf bytes.Buffer
			if str, ok := tc.body.(string); ok {
				buf.WriteString(str)
			} else {
				json.NewEncoder(&buf).Encode(tc.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/register", &buf)
			rec := httptest.NewRecorder()

			h.Register(rec, req)

			require.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedBody != nil {
				var res map[string]interface{}
				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedBody, res)
			}
		})
	}
}

func TestHandlers_Login(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		mockService    *mockAuthService
		expectedStatus int
		expectedCookie string
	}{
		{
			name: "Success",
			body: auth.LoginRequest{Email: "test@example.com", Password: "pass"},
			mockService: &mockAuthService{
				loginFunc: func(ctx context.Context, email, password string) (string, error) {
					return "valid-token", nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedCookie: "valid-token",
		},
		{
			name: "Invalid Body",
			body: "invalid-json",
			mockService: &mockAuthService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing Password",
			body: auth.LoginRequest{Email: "test@example.com"},
			mockService: &mockAuthService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid Credentials",
			body: auth.LoginRequest{Email: "test@example.com", Password: "wrong"},
			mockService: &mockAuthService{
				loginFunc: func(ctx context.Context, email, password string) (string, error) {
					return "", auth.ErrInvalidCredentials
				},
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Internal Error",
			body: auth.LoginRequest{Email: "test@example.com", Password: "pass"},
			mockService: &mockAuthService{
				loginFunc: func(ctx context.Context, email, password string) (string, error) {
					return "", errors.New("some error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := auth.NewHandlers(tc.mockService)

			var buf bytes.Buffer
			if str, ok := tc.body.(string); ok {
				buf.WriteString(str)
			} else {
				json.NewEncoder(&buf).Encode(tc.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/login", &buf)
			rec := httptest.NewRecorder()

			h.Login(rec, req)

			require.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedCookie != "" {
				cookies := rec.Result().Cookies()
				require.Len(t, cookies, 1)
				assert.Equal(t, "jwt", cookies[0].Name)
				assert.Equal(t, tc.expectedCookie, cookies[0].Value)
				assert.True(t, cookies[0].HttpOnly)
			}
		})
	}
}
