/*
* account é um aggregate root responsável por criar e controlar contas
 */

package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
	balance     Money
	currency    string
	active      bool
	isDefault   bool
	createdAt   time.Time
	updatedAt   time.Time
	events      []DomainEvent
}

// NewAccount cria uma nova conta
func NewAccount(userID uuid.UUID, name string, accountType AccountType, initialBalance Money) (*Account, error) {

	if userID == uuid.Nil {
		return nil, NewValidationError("user_id", "required")
	}

	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return nil, NewValidationError("name", "must be between 2 and 100 charactes")
	}

	validTypes := map[AccountType]bool{
		AccountTypeChecking:    true,
		AccountTypeSavings:     true,
		AccountTypeWallet:      true,
		AccountTypeCreditCard:  true,
		AccountTypeInvestiment: true,
	}

	if !validTypes[accountType] {
		return nil, NewValidationError("type", "invalid account type")
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

func (a *Account) UnsetDefault() {
	a.isDefault = false
}

func (a *Account) ClearEvents() {
	a.events = []DomainEvent{}
}

// Credit aumenta o saldo da conta (receita)
func (a *Account) Credit(amount Money, description string, categoryID *uuid.UUID, occurred_on time.Time) (*Transaction, error) {

	if err := a.checkActive(); err != nil {
		return nil, err
	}

	if err := a.checkCurrency(amount); err != nil {
		return nil, err
	}

	if amount.IsNegative() || amount.IsZero() {
		return nil, NewValidationError("amount", "must be positive for credit")
	}

	oldBalance := a.balance
	newBalance, err := a.balance.Add(amount)
	if err != nil {
		return nil, fmt.Errorf("Account.Credit: %w", err)
	}

	tx := NewTransaction(a.id, amount, TransactionTypeIncome, description, categoryID, occurred_on)
	a.balance = newBalance
	a.updatedAt = time.Now().UTC()

	a.addEvent(TransactionCreated{
		TransactionID: tx.ID(),
		AccountID:     a.id,
		Amount:        amount,
		Type:          TransactionTypeExpense,
		CategoryID:    categoryID,
		Description:   description,
		OccurredOn:    occurred_on,
	})

	a.addEvent(AccountBalanceUpdated{
		AccountID:  a.id,
		OldBalance: oldBalance,
		NewBalance: newBalance,
		OccurredOn: time.Now().UTC(),
	})

	return tx, nil

}

func (a *Account) Debit(amount Money, description string, categoryID *uuid.UUID, occurredOn time.Time) (*Transaction, error) {
	if err := a.checkActive(); err != nil {
		return nil, err
	}

	if err := a.checkCurrency(amount); err != nil {
		return nil, err
	}

	if amount.IsNegative() || amount.IsZero() {
		return nil, NewValidationError("amount", "must be positive for debit")
	}

	// Conta corrente e poupança não permitem saldo positivo
	if a.accountType == AccountTypeChecking || a.accountType == AccountTypeSavings {
		available := a.balance

		if amount.GreaterThan(available) {
			return nil, fmt.Errorf("account %s: %w", a.id, ErrInsufficientFunds)
		}
	}

	oldBalance := a.balance
	newBalance, err := a.balance.Sub(amount)
	if err != nil {
		return nil, fmt.Errorf("Account.Debit: %w", err)
	}

	tx := NewTransaction(a.id, amount.Abs(), TransactionTypeExpense, description, categoryID, occurredOn)

	//despesa é armazenada como negativo
	negAmount, err := NewMoney(-amount.Cents(), amount.Currency())
	if err != nil {
		return nil, err
	}
	tx.SetAmount(negAmount)

	a.balance = newBalance
	a.updatedAt = time.Now().UTC()

	a.addEvent(TransactionCreated{
		TransactionID: tx.ID(),
		AccountID:     a.id,
		Amount:        amount,
		Type:          TransactionTypeExpense,
		CategoryID:    categoryID,
		Description:   description,
		OccurredOn:    occurredOn,
	})

	a.addEvent(AccountBalanceUpdated{
		AccountID:  a.id,
		OldBalance: oldBalance,
		NewBalance: newBalance,
		OccurredOn: time.Now().UTC(),
	})

	return tx, nil

}

func (a *Account) Rename(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return NewValidationError("name", "must be between 2 ans 100 characters")
	}

	a.name = name
	a.updatedAt = time.Now().UTC()
	return nil
}

func (a *Account) Deactivate() error {
	if !a.IsActive() {
		return fmt.Errorf("account already inactive: %w", ErrInvalidOperation)
	}

	a.active = false
	a.updatedAt = time.Now().UTC()

	return nil
}

func (a *Account) addEvent(e DomainEvent) {
	a.events = append(a.events, e)
}

func (a *Account) checkCurrency(m Money) error {
	if m.Currency() != a.currency {
		return NewValidationError("currency", fmt.Sprintf("account currency is %s, got %s", a.currency, m.Currency()))
	}

	return nil
}

func (a *Account) checkActive() error {
	if !a.active {
		return fmt.Errorf("account %s is inactive: %w", a.id, ErrInvalidOperation)
	}

	return nil
}

func (a *Account) SetDefaultAccount() {
	a.isDefault = true
}

func (a *Account) ID() uuid.UUID            { return a.id }
func (a *Account) UserID() uuid.UUID        { return a.userID }
func (a *Account) Name() string             { return a.name }
func (a *Account) AccountType() AccountType { return a.accountType }
func (a *Account) Balance() Money           { return a.balance }
func (a *Account) Currency() string         { return a.currency }
func (a *Account) IsActive() bool           { return a.active }
func (a *Account) IsDefault() bool          { return a.isDefault }
func (a *Account) CreatedAt() time.Time     { return a.createdAt }
func (a *Account) UpdatedAt() time.Time     { return a.updatedAt }
func (a *Account) Events() []DomainEvent    { return a.events }

// ReconstructAccount reconstrói um Account a partir de dados persistidos.
// NÃO executa validações — use apenas ao carregar do banco.
// Isso é o padrão "reconstitution" do DDD: recriar o aggregate sem passar pelas regras de criação (que já foram validadas antes).
func ReconstructAccount(id uuid.UUID, userID uuid.UUID, name string, accountType AccountType, balance Money, isDefault bool, active bool, createdAt time.Time, updatedAt time.Time) *Account {
	return &Account{
		id:          id,
		userID:      userID,
		name:        name,
		accountType: accountType,
		balance:     balance,
		currency:    balance.Currency(),
		isDefault:   isDefault,
		active:      active,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
