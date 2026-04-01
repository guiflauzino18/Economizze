/*
* account é um aggregate root responsável por criar e controlar contas
 */

package aggregates

import (
	"time"

	"github.com/google/uuid"
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

}
