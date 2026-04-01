/*
* account é um aggregate root responsável por criar e controlar contas
 */

package aggregates

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/entities"
	"github.com/guiflauzino18/economizze/internal/domain/errors"
	"github.com/guiflauzino18/economizze/internal/domain/events"
	"github.com/guiflauzino18/economizze/internal/domain/vos"
)

type AccountType string

const (
	AccountTypeChecking    AccountType = "checking" // conta corrente
	AccountTypeSavings     AccountType = "savings"  // poupança
	AccountTypeWallet      AccountType = "wallet"   // carteira
	AccountTypeCreditCard  AccountType = "credit_card"
	AccountTypeInvestiment AccountType = "investment"
)

type Account struct {
	id          uuid.UUID
	userID      uuid.UUID
	name        string
	accountType AccountType
	balance     vos.Money
	currency    string
	active      bool
	createdAt   time.Time
	updatedAt   time.Time
	events      []events.DomainEvent
}

// NewAccount cria uma nova conta
func NewAccount(userID uuid.UUID, name string, accountType AccountType, initialBalance vos.Money) (*Account, error) {

	if userID == uuid.Nil {
		return nil, errors.NewValidationError("user_id", "required")
	}

	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return nil, errors.NewValidationError("name", "must be between 2 and 100 charactes")
	}

	validTypes := map[AccountType]bool{
		AccountTypeChecking:    true,
		AccountTypeSavings:     true,
		AccountTypeWallet:      true,
		AccountTypeCreditCard:  true,
		AccountTypeInvestiment: true,
	}

	if !validTypes[accountType] {
		return nil, errors.NewValidationError("type", "invalid account type")
	}

	now := time.Now().UTC()

	return &Account{
		id:          uuid.New(),
		userID:      userID,
		name:        name,
		accountType: accountType,
		balance:     initialBalance,
		currency:    initialBalance.Currency(),
		active:      true,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Credit aumenta o saldo da conta (receita)
func (a *Account) Credit() (*entities.Transaction, error) {

	return nil, nil
}
