package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain"
)

type CategoryRepository interface {
	FindAvailableForUser(ctx context.Context, userID uuid.UUID) ([]*domain.Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	Save(ctx context.Context, cat *domain.Category) error
}
