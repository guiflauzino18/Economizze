package ports

import (
	"context"

	"github.com/guiflauzino18/economizze/internal/domain"
)

type EventPublisher interface {
	Publish(ctx context.Context, events ...domain.DomainEvent) error
}
