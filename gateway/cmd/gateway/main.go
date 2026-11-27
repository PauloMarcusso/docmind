// Command gateway is the edge service of docmind and the composition root:
// it wires the adapters together and holds no business logic.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PauloMarcusso/docmind/gateway/internal/adapters/inbound/httpapi"
	"github.com/PauloMarcusso/docmind/gateway/internal/adapters/outbound/memory"
	"github.com/PauloMarcusso/docmind/gateway/internal/app"
	"github.com/PauloMarcusso/docmind/gateway/internal/domain"
)

// main only runs the gateway and turns a failure into exit code 1.
// os.Exit skips deferred calls, so it stays here and never inside run.
func main() {
	if err := run(context.Background(), os.Getenv); err != nil {
		slog.Error("gateway failed", "err", err)
		os.Exit(1)
	}
}

// run loads the config, wires the adapters and serves HTTP until the process
// gets SIGINT or SIGTERM, then shuts down gracefully within cfg.ShutdownTimeout.
func run(ctx context.Context, getenv func(string) string) error {
	cfg, err := loadConfig(getenv)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))

	// Wiring goes from the outside in: outbound adapter -> use case -> inbound adapter.
	store := memory.NewKeyStore(map[string]domain.Client{
		cfg.APIKey: {ID: "dev"},
	})
	auth := app.NewAuthenticator(store)
	router := httpapi.NewRouter(auth)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// The buffer lets the goroutine send and exit even if nobody reads errCh anymore.
	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	slog.Info("gateway listening", "addr", cfg.Addr)
	select {
	case err := <-errCh:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
	}
	slog.Info("gateway stopped")
	return nil
}
