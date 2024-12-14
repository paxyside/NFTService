package contract

import (
	"context"
	"log/slog"
	"time"
)

func (m *NFTContract) StartCacheUpdater(ctx context.Context, interval time.Duration) {
	l := slog.Default()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := m.updateTotalSupplyCache()
			if err != nil {
				l.Error("failed to update total supply cache", slog.Any("error", err))
			}
		case <-ctx.Done():
			l.Info("cache updater stopped")
			return
		}
	}
}
