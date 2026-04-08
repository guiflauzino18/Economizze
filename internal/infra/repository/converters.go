/*
* Conversores Dominio <-> Model
 */

package repository

import (
	"fmt"

	"github.com/guiflauzino18/economizze/internal/domain/aggregates"
	"github.com/guiflauzino18/economizze/internal/domain/entities"
	"github.com/guiflauzino18/economizze/internal/domain/vos"
)

// ============================================================
// Account converters
// ============================================================

// accountToModel converte a entidade de domínio para o modelo GORM.
// Chamado antes de salvar no banco.
func accountToModel(a aggregates.Account) AccountModel {
	return AccountModel{
		ID:           a.ID(),
		UserID:       a.UserID(),
		Name:         a.Name(),
		AccountType:  string(a.AccountType()),
		BalanceCents: a.Balance().Cents(),
		Currency:     a.Balance().Currency(),
		IsDefault:    a.IsDefault(),
		Active:       a.IsActive(),
		CreatedAt:    a.CreatedAt(),
		UpdatedAt:    a.UpdatedAt(),
	}
}

// modelToAccount converte o modelo GORM de volta para a entidade de domínio.
// Chamado após buscar do banco.
// Usa ReconstructAccount para reconstruir sem passar pelas validações
// do construtor (o dado já foi validado quando foi salvo).
func modelToAccount(m AccountModel) (*aggregates.Account, error) {
	balance, err := vos.NewMoney(m.BalanceCents, m.Currency)
	if err != nil {
		return nil, fmt.Errorf("modelToAccount.balance: %w", err)
	}

	// ReconstructAccount bypassa as validações do NewAccount
	// porque estamos reconstruindo um dado já válido do banco
	return aggregates.ReconstructAccount(
		m.ID,
		m.UserID,
		m.Name,
		aggregates.AccountType(m.AccountType),
		balance,
		m.IsDefault,
		m.Active,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

// ============================================================
// Transaction converters
// ============================================================

func transactionToModel(t *entities.Transaction) TransactionModel {
	return TransactionModel{
		ID:             t.ID(),
		AccountID:      t.AccountID(),
		CategoryID:     t.CategoryID(),
		TransferPeerID: t.TransferPeerID(),
		Cents:          t.Amount().Cents(),
		Currency:       t.Amount().Currency(),
		Type:           string(t.Type()),
		Description:    t.Description(),
		Notes:          t.NotesPtr(), // retorna *string
		OccurredOn:     t.OccurredOn(),
		RecurringID:    t.RecurringID(),
		CreatedAt:      t.CreatedAt(),
		UpdatedAt:      t.UpdatedAt(),
	}
}

func modelToTransaction(m TransactionModel) (*entities.Transaction, error) {
	amount, err := vos.NewMoney(m.Cents, m.Currency)
	if err != nil {
		return nil, fmt.Errorf("modelToTransaction.amount: %w", err)
	}

	return entities.ReconstructTransaction(
		m.ID,
		m.AccountID,
		m.CategoryID,
		m.TransferPeerID,
		amount,
		entities.TransactionType(m.Type),
		m.Description,
		m.Notes,
		m.OccurredOn,
		m.RecurringID,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

// ============================================================
// Budget converters
// ============================================================

func budgetToModel(b *aggregates.Budget) BudgetModel {
	return BudgetModel{
		ID:                 b.ID(),
		UserID:             b.UserID(),
		CategoryID:         b.CategoryID(),
		PeriodStart:        b.Period().Start(),
		PeriodEnd:          b.Period().End(),
		LimitCents:         b.Limit().Cents(),
		SpentCents:         b.Spent().Cents(),
		Currency:           b.Limit().Currency(),
		NotifyWhenExceeded: b.NotifyWhenExceeded(),
		CreatedAt:          b.CreatedAt(),
		UpdatedAt:          b.UpdatedAt(),
	}
}

func modelToBudget(m BudgetModel) (*aggregates.Budget, error) {
	limit, err := vos.NewMoney(m.LimitCents, m.Currency)
	if err != nil {
		return nil, fmt.Errorf("modelToBudget.limit: %w", err)
	}

	spent, err := vos.NewMoney(m.SpentCents, m.Currency)
	if err != nil {
		return nil, fmt.Errorf("modelToBudget.spent: %w", err)
	}

	period, err := vos.NewPediod(m.PeriodStart, m.PeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("modelToBudget.period: %w", err)
	}

	return aggregates.ReconstructBudget(
		m.ID,
		m.UserID,
		m.CategoryID,
		period,
		limit,
		spent,
		m.NotifyWhenExceeded,
		m.CreatedAt,
		m.UpdatedAt,
	), nil

}

// ============================================================
// Category converters
// ============================================================

func CategoryToModel(c *entities.Category) CategoryModel {
	return CategoryModel{
		ID:          c.ID(),
		UserID:      c.UserID(),
		Name:        c.Name(),
		DefaultType: string(c.DefaultType()),
		Active:      c.IsActive(),
		CreatedAt:   c.CreatedAt(),
		UpdatedAt:   c.UpdatedAt(),
	}
}

func modelToCategory(c CategoryModel) *entities.Category {
	return entities.ReconstructCategory(
		c.ID,
		c.UserID,
		c.Name,
		entities.TransactionType(c.DefaultType),
		c.Active,
		c.CreatedAt,
		c.UpdatedAt,
	)
}
