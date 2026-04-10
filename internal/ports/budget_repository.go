package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain"
)

type BudgetRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Budget, error)
	FindByUserAndPeriod(ctx context.Context, userID uuid.UUID, period domain.Period) ([]*domain.Budget, error)
	Save(ctx context.Context, budget *domain.Budget) error
	UpdateSpent(ctx context.Context, budgetID uuid.UUID, spentCents int64) error
}
