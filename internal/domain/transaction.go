/*
* Transaction é a entidade que representa cada transação realizada
 */

package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

type Transaction struct {
	id             uuid.UUID
	accountID      uuid.UUID
	amount         Money // positivo = receita, negativo = despesa
	txType         TransactionType
	description    string
	categoryID     *uuid.UUID
	transferPeerID *uuid.UUID
	recurringID    *uuid.UUID
	occurredOn     time.Time
	createdAt      time.Time
	updatedAt      time.Time
	notes          string
}

// NewTransaction cria uma nova transação
func NewTransaction(accountID uuid.UUID, amount Money, txType TransactionType, description string, categoryID *uuid.UUID, occurredOn time.Time) *Transaction {
	return &Transaction{
		id:          uuid.New(),
		accountID:   accountID,
		amount:      amount,
		txType:      txType,
		description: description,
		categoryID:  categoryID,
		occurredOn:  occurredOn,
		createdAt:   time.Now().UTC(),
	}
}

// AddNote adiciona nota
func (t *Transaction) AddNote(note string) {
	t.notes = note
}

// Categorize altera categoria
func (t *Transaction) Categorize(categoryID uuid.UUID) {
	t.categoryID = &categoryID
}

// getters
func (t *Transaction) ID() uuid.UUID              { return t.id }
func (t *Transaction) AccountID() uuid.UUID       { return t.accountID }
func (t *Transaction) Amount() Money              { return t.amount }
func (t *Transaction) Type() TransactionType      { return t.txType }
func (t *Transaction) Description() string        { return t.description }
func (t *Transaction) CategoryID() *uuid.UUID     { return t.categoryID }
func (t *Transaction) TransferPeerID() *uuid.UUID { return t.transferPeerID }
func (t *Transaction) OccurredOn() time.Time      { return t.occurredOn }
func (t *Transaction) RecurringID() *uuid.UUID    { return t.recurringID }
func (t *Transaction) CreatedAt() time.Time       { return t.createdAt }
func (t *Transaction) UpdatedAt() time.Time       { return t.updatedAt }
func (t *Transaction) NotesPtr() *string          { return &t.notes }

// setters
func (t *Transaction) SetAmount(a Money) {
	t.amount = a
}

func ReconstructTransaction(id uuid.UUID, accountID uuid.UUID, categoryID *uuid.UUID, transferPeerID *uuid.UUID, amount Money, txType TransactionType, description string, notes *string, occurredOn time.Time, recurringID *uuid.UUID, createdAt time.Time, updatedAt time.Time) *Transaction {
	t := &Transaction{
		id:             id,
		accountID:      accountID,
		categoryID:     categoryID,
		transferPeerID: transferPeerID,
		amount:         amount,
		txType:         txType,
		description:    description,
		occurredOn:     occurredOn,
		recurringID:    recurringID,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
	if notes != nil {
		t.notes = *notes
	}
	return t
}
