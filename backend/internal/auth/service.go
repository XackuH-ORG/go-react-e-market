package auth

import (
	"context"
	"errors"
	"time"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
	"github.com/XackuH-ORG/go-react-e-market/backend/internal/utils"
)

type Service struct {
	db *database.Queries
}

func NewService(db *database.Queries) *Service {
	return &Service{db: db}
}

func (s *Service) Register(ctx context.Context, email, password string) (*database.User, error) {
	// 1. Хэшируем пароль
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 2. Создаем юзера с дефолтной ролью CUSTOMER
	// Предполагается, что в БД тип user_role это ENUM. Если это строка, то просто "CUSTOMER"
	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		Role:         database.UserRoleCUSTOMER, // Убедись, что константа совпадает с тем, что сгенерировал sqlc
	})
	if err != nil {
		return nil, err // Здесь можно добавить проверку на unique constraint (уже существует)
	}

	return &user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	// 1. Ищем юзера
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("неверный email или пароль")
	}

	// 2. Проверяем пароль
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("неверный email или пароль")
	}

	// 3. Генерируем токен
	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return "", err
	}

	// 4. Сохраняем сессию (например, на 30 дней)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	// Важно: если sqlc сгенерировал pgtype.Timestamp, нужно будет обернуть время
	// Если ты используешь обычный time.Time в sqlc (через overrides), то оставляем так
	_, err = s.db.CreateSession(ctx, database.CreateSessionParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
