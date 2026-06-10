package products

import (
	repo "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type Service interface {
}

type svc struct {
	repo repo.Querier
}

func NewService(db repo.Querier) Service {
	return &svc{repo: db}
}
