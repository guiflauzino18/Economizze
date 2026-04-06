package aggregates

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/errors"
	"github.com/guiflauzino18/economizze/internal/domain/events"
	"github.com/guiflauzino18/economizze/internal/domain/vos"
)

type Budget struct {
	id         uuid.UUID
	userID     uuid.UUID
	categoryID uuid.UUID
	period     vos.Period
	limit      vos.Money
	spent      vos.Money
	events     []events.DomainEvent
	createdAt  time.Time
	updatedAt  time.Time
}

func NewBudget(userID, categoryID uuid.UUID, period vos.Period, limit vos.Money) (*Budget, error) {
	if userID == uuid.Nil {
		return nil, errors.NewValidationError("user_id", "required")
	}

	if categoryID == uuid.Nil {
		return nil, errors.NewValidationError("category_id", "required")
	}

	if !limit.IsPositive() {
		return nil, errors.NewValidationError("limit", "must be positive")
	}

	zero, _ := vos.NewMoney(0, limit.Currency())
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
func (b *Budget) RegisterSpending(amount vos.Money) error {
	if amount.Currency() != b.limit.Currency() {
		return errors.NewValidationError("currency", "currency mismatch")
	}

	if !amount.IsPositive() {
		return errors.NewValidationError("amount", "must be positive")
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
		b.addEvent(events.BudgetExceeded{
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

func (b *Budget) ID() uuid.UUID                { return b.id }
func (b *Budget) UserID() uuid.UUID            { return b.userID }
func (b *Budget) CategoryID() uuid.UUID        { return b.categoryID }
func (b *Budget) Period() vos.Period           { return b.period }
func (b *Budget) Limit() vos.Money             { return b.limit }
func (b *Budget) Spent() vos.Money             { return b.spent }
func (b *Budget) Events() []events.DomainEvent { return b.events }
func (b *Budget) ClearEvents()                 { b.events = nil }

func (b *Budget) addEvent(e events.DomainEvent) {
	b.events = append(b.events, e)
}
