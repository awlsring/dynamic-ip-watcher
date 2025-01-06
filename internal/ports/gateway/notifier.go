package gateway

import (
	"context"

	"github.com/awlsring/dynamic-ip-watcher/internal/core/domain/event"
)

type Notifier interface {
	SendEventMessage(ctx context.Context, event event.Event) error
}
