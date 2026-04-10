package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain"
)

type TransactionFilter struct {
	AccountID  *uuid.UUID
	CategoryID *uuid.UUID
	Type       *domain.TransactionType
	From       *time.Time
	To         *time.Time
	Search     string
	Page, Size int
}

type TransactionRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	FindAll(ctx context.Context, filter TransactionFilter) ([]*domain.Transaction, int64, error)
	Save(ctx context.Context, tx *domain.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
}
