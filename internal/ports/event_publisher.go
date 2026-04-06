package ports

import (
	"context"

	"github.com/guiflauzino18/economizze/internal/domain/events"
)

type EventPublisher interface {
	Publish(ctx context.Context, events ...events.DomainEvent) error
}
