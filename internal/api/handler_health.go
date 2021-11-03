package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

// HealthzHandler is a k8s liveness endpoint
func HealthzHandler(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, HealthStatusResponse{Status: StatusHealthy})
}

// ReadyzHandler is a k8s readiness endpoint
func ReadyzHandler(fns ...func(ctx context.Context) HealthEntryStatus) func(echo.Context) error {
	return func(ctx echo.Context) error {
		entries := make([]HealthEntryStatus, 0, len(fns))
		ch := make(chan HealthEntryStatus, len(fns))
		var wg sync.WaitGroup
		healthy := true

		for _, fn := range fns {
			wg.Add(1)
			go func(f func(context.Context) HealthEntryStatus, w *sync.WaitGroup) {
				defer w.Done()
				ch <- f(ctx.Request().Context())
			}(fn, &wg)
		}
		wg.Wait()
		close(ch)

		for entry := range ch {
			healthy = healthy && entry.Status == StatusHealthy
			entries = append(entries, entry)
		}

		if healthy {
			return ctx.JSON(http.StatusOK, HealthStatusResponse{
				Status:  StatusHealthy,
				Entries: entries,
			})
		}
		return ctx.JSON(http.StatusServiceUnavailable, HealthStatusResponse{
			Status:  StatusUnhealthy,
			Entries: entries,
		})
	}

}
