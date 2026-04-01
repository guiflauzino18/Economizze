/*
* Transaction é a entidade que representa cada transação realizada
 */

package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/vos"
)

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

type Transaction struct {
	id          uuid.UUID
	accountID   uuid.UUID
	amount      vos.Money // positivo = receita, negativo = despesa
	txType      TransactionType
	description string
	categoryID  *uuid.UUID
	occurredOn  time.Time
	createdAt   time.Time
	notes       string
}

// NewTransaction cria uma nova transação
func NewTransaction(accountID uuid.UUID, amount vos.Money, txType TransactionType, description string, categoryID *uuid.UUID, occurredOn time.Time) *Transaction {
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
func (t *Transaction) ID() uuid.UUID          { return t.id }
func (t *Transaction) AccountID() uuid.UUID   { return t.accountID }
func (t *Transaction) Amount() vos.Money      { return t.amount }
func (t *Transaction) Type() TransactionType  { return t.txType }
func (t *Transaction) Description() string    { return t.description }
func (t *Transaction) CategoryID() *uuid.UUID { return t.categoryID }
func (t *Transaction) OccurredOn() time.Time  { return t.occurredOn }
func (t *Transaction) Notes() string          { return t.notes }
