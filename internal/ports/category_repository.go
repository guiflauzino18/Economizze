package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/entities"
)

type CategoryRepository interface {
	FindAvailableForUser(ctx context.Context, userID uuid.UUID) ([]*entities.Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Category, error)
	Save(ctx context.Context, cat *entities.Category) error
}
