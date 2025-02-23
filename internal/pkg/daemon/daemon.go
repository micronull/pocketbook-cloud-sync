package daemon

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/command/sync/factory"
)

type Daemon struct {
	timeout time.Duration
	sync    factory.Synchronizer
}

var _ factory.Synchronizer = (*Daemon)(nil)

// New construct for [Daemon].
// timeout - how long to wait between synchronizations.
func New(timeout time.Duration, sync factory.Synchronizer) *Daemon {
	return &Daemon{
		timeout: timeout,
		sync:    sync,
	}
}

func (d Daemon) Sync(ctx context.Context) error {
	ticker := time.NewTicker(d.timeout)
	defer ticker.Stop()

	for {
		if err := d.sync.Sync(ctx); err != nil {
			var httpErr interface {
				Code() int
			}

			if errors.As(err, &httpErr) && httpErr.Code() >= http.StatusInternalServerError {
				slog.Error("server error", "error", err)
			} else {
				return fmt.Errorf("call sync: %w", err)
			}
		}

		ticker.Reset(d.timeout)

		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}
