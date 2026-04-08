package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/entities"
	domainErrors "github.com/guiflauzino18/economizze/internal/domain/errors"
	"github.com/guiflauzino18/economizze/internal/ports"
	"gorm.io/gorm"
)

var _ ports.TransactionRepository = (*transactionRepository)(nil)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ports.TransactionRepository {
	return &transactionRepository{db}
}

// Delete implements [ports.TransactionRepository].
func (t *transactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := t.db.WithContext(ctx).Delete(&TransactionModel{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("transactionRepository.Delete %s: %w", id, domainErrors.ErrNotFound)
	}

	return result.Error
}

// FindAll implements [ports.TransactionRepository].
func (t *transactionRepository) FindAll(ctx context.Context, filter ports.TransactionFilter) ([]*entities.Transaction, int64, error) {
	var models []TransactionModel

	var total int64

	// Constrói a query base com filtros opcionais
	q := t.db.WithContext(ctx).Model(&TransactionModel{})

	if filter.AccountID != nil {
		q = q.Where("account_id = ?", filter.AccountID)
	}

	if filter.CategoryID != nil {
		q = q.Where("category_id = ?", filter.AccountID)
	}

	if filter.Type != nil {
		q = q.Where("type = ?", filter.Type)
	}

	if filter.From != nil {
		q = q.Where("occurred_on = ?", filter.From)
	}

	if filter.To != nil {
		q = q.Where("occurred_on <= ?", filter.To)
	}

	if filter.Search != "" {
		// ILIKE com trigram index para busca parcial eficiente
		// O índice idx_transactions_description_trgm criado na migration 000008 garante performance mesmo em tabelas grandes
		q = q.Where("descritprion ILIKE ?", "%"+filter.Search+"%")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("transactionRepository.FindAll.count: %w", err)
	}

	// Aplica paginação e ordenação
	page, size := normalizePaginator(filter.Page, filter.Size)

	err := q.
		Preload("Category").
		Order("occurred_on DESC, created_at DESC").
		Limit(size).
		Offset(page - 1*size).
		Find(&models).Error

	if err != nil {
		return nil, 0, fmt.Errorf("transactionRepository.FindAll.find: %w", err)
	}

	transactions := make([]*entities.Transaction, 0, len(models))

	for _, m := range models {
		tx, err := modelToTransaction(m)
		if err != nil {
			return nil, 0, fmt.Errorf("transactionRepository.FindAll.convert: %w", err)
		}

		transactions = append(transactions, tx)
	}

	return transactions, total, nil

}

// FindByID implements [ports.TransactionRepository].
func (t *transactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Transaction, error) {
	var model TransactionModel

	err := t.db.WithContext(ctx).
		Preload("Category").
		First(&model, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("transactionRepository.FindByID %s: %w", id, domainErrors.ErrNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("transactionRepository.FindByID %w", err)
	}

	return modelToTransaction(model)
}

// Save implements [ports.TransactionRepository].
func (t *transactionRepository) Save(ctx context.Context, tx *entities.Transaction) error {
	model := transactionToModel(tx)

	err := t.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return fmt.Errorf("transactionRepository.save %s: %w", tx.ID(), err)
	}

	return nil
}
