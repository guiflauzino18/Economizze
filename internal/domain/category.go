package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	id          uuid.UUID
	userID      uuid.UUID
	name        string
	defaultType TransactionType
	active      bool
	createdAt   time.Time
	updatedAt   time.Time
}

func NewCategory(usreID uuid.UUID, name string, typeTx TransactionType) (*Category, error) {

	if len(name) < 2 {
		return nil, NewValidationError("name", "must be at least 2 characters")
	}

	return &Category{
		id:          uuid.New(),
		userID:      usreID,
		name:        name,
		defaultType: typeTx,
		active:      true,
		createdAt:   time.Now().UTC(),
		updatedAt:   time.Now().UTC(),
	}, nil
}

func (c *Category) ID() uuid.UUID                { return c.id }
func (c *Category) UserID() uuid.UUID            { return c.userID }
func (c *Category) Name() string                 { return c.name }
func (c *Category) DefaultType() TransactionType { return c.defaultType }
func (c *Category) IsActive() bool               { return c.active }
func (c *Category) CreatedAt() time.Time         { return c.createdAt }
func (c *Category) UpdatedAt() time.Time         { return c.updatedAt }

func ReconstructCategory(id uuid.UUID, userID uuid.UUID, name string, defaultType TransactionType, active bool, createdAt time.Time, updatedAt time.Time) *Category {
	return &Category{
		id:          id,
		userID:      userID,
		name:        name,
		defaultType: defaultType,
		active:      active,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
