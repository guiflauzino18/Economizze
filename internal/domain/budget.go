package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Budget struct {
	id                 uuid.UUID
	userID             uuid.UUID
	categoryID         uuid.UUID
	period             Period
	limit              Money
	spent              Money
	notifyWhenExceeded bool
	events             []DomainEvent
	createdAt          time.Time
	updatedAt          time.Time
}

func NewBudget(userID, categoryID uuid.UUID, period Period, limit Money) (*Budget, error) {
	if userID == uuid.Nil {
		return nil, NewValidationError("user_id", "required")
	}

	if categoryID == uuid.Nil {
		return nil, NewValidationError("category_id", "required")
	}

	if !limit.IsPositive() {
		return nil, NewValidationError("limit", "must be positive")
	}

	zero, _ := NewMoney(0, limit.Currency())
	now := time.Now().UTC()

	return &Budget{
		id:         uuid.New(),
		userID:     userID,
		categoryID: categoryID,
		period:     period,
		limit:      limit,
		spent:      zero,
		createdAt:  now,
		updatedAt:  now,
	}, nil

}

// RegisterSpending - Registra um gasto e verifica se o budget foi excedido
func (b *Budget) RegisterSpending(amount Money) error {
	if amount.Currency() != b.limit.Currency() {
		return NewValidationError("currency", "currency mismatch")
	}

	if !amount.IsPositive() {
		return NewValidationError("amount", "must be positive")
	}

	newSpent, err := b.spent.Add(amount)
	if err != nil {
		return fmt.Errorf("Budget.RegisterSpending: %w", err)
	}

	wasUnder := !b.spent.GreaterThan(b.limit)
	b.spent = newSpent
	b.updatedAt = time.Now().UTC()

	// Dispara evento se budget exceeded
	if wasUnder && b.spent.GreaterThan(b.limit) {
		b.addEvent(BudgetExceeded{
			BudgetID:   b.id,
			CategoryID: b.categoryID,
			Limit:      b.limit,
			Spent:      b.spent,
			OccurredOn: time.Now().UTC(),
		})
	}

	return nil

}

// PercentUsed - 0.0 a 1.0+ (pode ser mais de 1 se excedido)
func (b *Budget) PercentUsed() float64 {
	if b.limit.Cents() == 0 {
		return 0
	}

	return float64(b.spent.Cents()) / float64(b.limit.Cents())
}

func (b *Budget) ID() uuid.UUID            { return b.id }
func (b *Budget) UserID() uuid.UUID        { return b.userID }
func (b *Budget) CategoryID() uuid.UUID    { return b.categoryID }
func (b *Budget) Period() Period           { return b.period }
func (b *Budget) Limit() Money             { return b.limit }
func (b *Budget) Spent() Money             { return b.spent }
func (b *Budget) NotifyWhenExceeded() bool { return b.notifyWhenExceeded }
func (b *Budget) Events() []DomainEvent    { return b.events }
func (b *Budget) ClearEvents()             { b.events = nil }
func (b *Budget) CreatedAt() time.Time     { return b.createdAt }
func (b *Budget) UpdatedAt() time.Time     { return b.updatedAt }

func (b *Budget) addEvent(e DomainEvent) {
	b.events = append(b.events, e)
}

func ReconstructBudget(
	id uuid.UUID,
	userID uuid.UUID,
	categoryID uuid.UUID,
	period Period,
	limit Money,
	spent Money,
	notifyWhenExceeded bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Budget {
	return &Budget{
		id:                 id,
		userID:             userID,
		categoryID:         categoryID,
		period:             period,
		limit:              limit,
		spent:              spent,
		notifyWhenExceeded: notifyWhenExceeded,
		createdAt:          createdAt,
		updatedAt:          updatedAt,
	}
}
