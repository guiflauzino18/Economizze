package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/aggregates"
)

type AccountRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*aggregates.Account, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*aggregates.Account, error)
	Save(ctx context.Context, account *aggregates.Account) error
	Delete(ctx context.Context, id uuid.UUID)
}
