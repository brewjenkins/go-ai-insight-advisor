package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"brewjenkins/go-ai-insight-advisor/internal/app"
	"brewjenkins/go-ai-insight-advisor/internal/config"
	"brewjenkins/go-ai-insight-advisor/internal/database"
	"brewjenkins/go-ai-insight-advisor/internal/domain/constant"

	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listener",
	Short: "Start the HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		logger := newLogger(cfg.App.Env)

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		startupCtx, cancelStartup := context.WithTimeout(ctx, 10*time.Second)
		defer cancelStartup()
		db, err := database.OpenPostgres(startupCtx, cfg.Postgres)
		if err != nil {
			return fmt.Errorf("open PostgreSQL: %w", err)
		}
		defer db.Close()

		router := app.NewRouter(cfg, db, logger)
		addr := net.JoinHostPort(cfg.App.Server.Host, cfg.App.Server.Port)

		srv := &http.Server{
			Addr:              addr,
			Handler:           router,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		}

		errCh := make(chan error, 1)
		go func() {
			logger.Info("http server starting", "addr", addr, "env", cfg.App.Env)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
			close(errCh)
		}()

		select {
		case err := <-errCh:
			if err != nil {
				return fmt.Errorf("http server: %w", err)
			}
		case <-ctx.Done():
			logger.Info("shutdown signal received, draining connections")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("graceful shutdown: %w", err)
			}
			logger.Info("http server stopped")
		}

		return nil
	},
}

func newLogger(env string) *slog.Logger {
	if env == constant.MODE_PRODUCTION {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
