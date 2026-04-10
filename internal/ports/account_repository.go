package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain"
)

type AccountRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error)
	Save(ctx context.Context, account *domain.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
}
