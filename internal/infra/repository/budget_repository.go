package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain"
	"github.com/guiflauzino18/economizze/internal/ports"
	"gorm.io/gorm"
)

var _ ports.BudgetRepository = (*budgetRepository)(nil)

type budgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) ports.BudgetRepository {
	return &budgetRepository{db}
}

// Delete implements [ports.BudgetRepository].
func (b *budgetRepository) UpdateSpent(ctx context.Context, budgetID uuid.UUID, spentCents int64) error {
	// Atualiza só o campo spent_cents — mais eficiente que Save completo
	// Chamado pelo use case após cada transação na categoria do orçamento
	result := b.db.WithContext(ctx).
		Model(&BudgetModel{}).
		Where("id = ?", budgetID).
		Updates(map[string]any{
			"spent_cents": spentCents,
			"updated_at":  time.Now().UTC(),
		})

	if result.RowsAffected == 0 {
		return fmt.Errorf("budgetRepository.UpdateSpent %s: %w", budgetID, domain.ErrNotFound)
	}

	return result.Error
}

// FindByID implements [ports.BudgetRepository].
func (b *budgetRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Budget, error) {
	var model BudgetModel

	err := b.db.WithContext(ctx).
		Preload("Category").
		First(&model, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("budgetRespository.FinByID %s: %w", id, domain.ErrNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("budgetRepository.FindByID: %w", err)
	}

	return modelToBudget(model)
}

// FindByUserAndPeriod implements [ports.BudgetRepository].
func (b *budgetRepository) FindByUserAndPeriod(ctx context.Context, userID uuid.UUID, period domain.Period) ([]*domain.Budget, error) {
	var models []BudgetModel

	err := b.db.WithContext(ctx).
		Preload("Category").
		Where(`
			user_id = ? AND
			period_start >= ? AND
			period_end <= ?
		`, userID, period.Start(), period.End()).
		Order("period_start ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("budgetRepository.FindByUserAndPeriod: %w", err)
	}

	budgets := make([]*domain.Budget, 0, len(models))

	for _, m := range models {
		budget, err := modelToBudget(m)
		if err != nil {
			return nil, fmt.Errorf("budgetRepository.FindByUserAndPeriod.convert: %w", err)
		}

		budgets = append(budgets, budget)
	}

	return budgets, nil
}

// Save implements [ports.BudgetRepository].
func (b *budgetRepository) Save(ctx context.Context, budget *domain.Budget) error {
	model := budgetToModel(budget)

	if err := b.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("budgetRepository.Save: %w", err)
	}

	return nil
}
