package products

import (
	"context"

	"github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) error
}

type svc struct {
	db repo.Querier
}

func NewService(db repo.Querier) Service {
	return &svc{
		db: db,
	}
}

func (s *svc) ListProducts(ctx context.Context) error {
	return nil
}

