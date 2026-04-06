package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/aggregates"
	"github.com/guiflauzino18/economizze/internal/domain/vos"
)

type BudgetRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*aggregates.Budget, error)
	FindByUserAndPeriod(ctx context.Context, userID uuid.UUID, period vos.Period)
	Save(ctx context.Context, budget *aggregates.Budget) error
	Delete(ctx context.Context, id uuid.UUID) error
}
