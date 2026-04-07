package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/aggregates"
	domainErrors "github.com/guiflauzino18/economizze/internal/domain/errors"
	"github.com/guiflauzino18/economizze/internal/ports"
	"gorm.io/gorm"
)

var _ ports.AccountRepository = (*accountRepository)(nil)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) ports.AccountRepository {
	return &accountRepository{db}
}

// Delete implements [ports.AccountRepository].
func (a *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := a.db.WithContext(ctx).
		Model(&AccountModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"active":     false,
			"updated_at": time.Now().UTC(),
		})

	if result.Error != nil {
		return fmt.Errorf("accountRepository.Delete %s: %w", id, result.Error)
	}

	// Se 0 então não encontrou o ID
	if result.RowsAffected == 0 {
		return fmt.Errorf("accountRepository.Delete %s: %w", id, domainErrors.ErrNotFound)
	}

	return nil
}

// FindByID implements [ports.AccountRepository].
func (a *accountRepository) FindByID(ctx context.Context, id uuid.UUID) (*aggregates.Account, error) {
	var model AccountModel

	err := a.db.WithContext(ctx).
		Where("id = ? AND active = true", id).
		First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("accountRepository.FindByID %s: %w", id, domainErrors.ErrNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("accountRepository.FindByID %s: %w", id, err)
	}

	account, err := modelToAccount(model)
	if err != nil {
		return nil, fmt.Errorf("accountRepository.FindByID.convert: %w", err)
	}

	return account, nil
}

// FindByUserID implements [ports.AccountRepository].
func (a *accountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*aggregates.Account, error) {
	var models []AccountModel

	err := a.db.WithContext(ctx).Where("user_id = ? AND active = true", userID).
		Order("is_default DESC, created_at ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("accountRepository.FindByUserID %s: %w", userID, err)
	}

	accounts := make([]*aggregates.Account, 0, len(models))

	for _, m := range models {
		a, err := modelToAccount(m)
		if err != nil {
			return nil, fmt.Errorf("accountRepository.FindByID.convert: %w", err)
		}

		accounts = append(accounts, a)
	}

	return accounts, nil

}

// Save implements [ports.AccountRepository].
func (a *accountRepository) Save(ctx context.Context, account *aggregates.Account) error {
	model := accountToModel(*account)

	// Save faz INSERT ou UPDATE automaticamente baseado na PK:
	// - PK zero → INSERT
	// - PK preenchida → UPDATE de todos os campos
	// Usamos Save em vez de Create para suportar upsert simples

	err := a.db.WithContext(ctx).Save(&model).Error

	if err != nil {
		return fmt.Errorf("accountRepository.Save %s: %w", account.ID(), err)
	}

	return nil
}
