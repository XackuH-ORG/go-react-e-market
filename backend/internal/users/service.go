package users

import (
	"context"
	"errors"

	"github.com/google/uuid"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type Service struct {
	db *database.Queries
}

func NewService(db *database.Queries) *Service {
	return &Service{db: db}
}

// ListUsers возвращает список пользователей с пагинацией
func (s *Service) ListUsers(ctx context.Context, limit, offset int32) ([]database.GetUsersRow, error) {
	return s.db.GetUsers(ctx, database.GetUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

// UpdateRole обновляет роль пользователя
func (s *Service) UpdateRole(ctx context.Context, userID uuid.UUID, role string) (*database.UpdateUserRoleRow, error) {
	// Валидация роли
	var userRole database.UserRole
	switch role {
	case string(database.UserRoleADMIN):
		userRole = database.UserRoleADMIN
	case string(database.UserRoleCUSTOMER):
		userRole = database.UserRoleCUSTOMER
	default:
		return nil, errors.New("недопустимая роль")
	}

	// Выполняем запрос
	user, err := s.db.UpdateUserRole(ctx, database.UpdateUserRoleParams{
		ID:   userID,
		Role: userRole,
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}
