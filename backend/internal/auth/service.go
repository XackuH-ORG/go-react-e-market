package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*repo.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	q         repo.Querier
	jwtSecret []byte
}

func NewAuthService(q repo.Querier, jwtSecret []byte) AuthService {
	return &authService{
		q:         q,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (*repo.User, error) {
	// Check if user already exists. We ignore the error here because sql.ErrNoRows is expected if user doesn't exist.
	_, err := s.q.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.q.CreateUser(ctx, repo.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
		Role:         repo.UserRoleCUSTOMER, // Default role
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		// Either user not found or DB error, treat as invalid credentials
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
