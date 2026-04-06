package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/entities"
)

type TransactionFilter struct {
	AccountID    *uuid.UUID
	CategoryID   *uuid.UUID
	Type         *entities.TransactionType
	From         *time.Time
	To           *time.Time
	Search, Size int
}

type TransactionRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Transaction, error)
	FindAll(ctx context.Context, filter TransactionFilter) ([]*entities.Transaction, int64, error)
	Save(ctx context.Context, tx *entities.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
}
