/*
* events disparam eventos quando realizado alguma ação.
 */

package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/vos"
)

// DomainEvent é um contrato base
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
	AggregateID() uuid.UUID
}

// TransactionCreated - dispara ao criar qualquer transação
type TransactionCreated struct {
	TransactionID uuid.UUID
	AccountID     uuid.UUID
	Amount        vos.Money
	CategoryID    *uuid.UUID
	Description   string
	OccurredOn    time.Time
}

func (e TransactionCreated) EventName() string {
	return "transaction.created"
}

func (e TransactionCreated) OccurredAt() time.Time {
	return e.OccurredOn
}

func (e TransactionCreated) AggregateID() uuid.UUID {
	return e.AccountID
}

// AccountBalanceUpdated - disparado quando o saldo muda
type AccountBalanceUpdated struct {
	AccountID  uuid.UUID
	OldBalance vos.Money
	NewBalance vos.Money
	OccurredOn time.Time
}

func (e AccountBalanceUpdated) EventName() string {
	return "account.balance_updated"
}

func (e AccountBalanceUpdated) OccurredAt() time.Time {
	return e.OccurredOn
}

func (e AccountBalanceUpdated) AggregateID() uuid.UUID {
	return e.AccountID
}

// BudgetExceeded - disparado quando extrapola o orçamento
type BudgetExceeded struct {
	BudgetID   uuid.UUID
	CategoryID uuid.UUID
	Limit      vos.Money
	Spent      vos.Money
	OccurredOn time.Time
}

func (e BudgetExceeded) EventName() string {
	return "budget.exceeded"
}

func (e BudgetExceeded) OccurredAt() time.Time {
	return e.OccurredOn
}

func (e BudgetExceeded) AggregateID() uuid.UUID {
	return e.BudgetID
}
