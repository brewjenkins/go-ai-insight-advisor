package database

import (
	"context"
	"errors"
	"strings"
	"testing"

	"brewjenkins/go-ai-insight-advisor/internal/config"
)

func TestOpenPostgresInvalidSSLMode(t *testing.T) {
	cfg := config.NewConfig().Postgres
	cfg.SSLMode = "invalid"
	cfg.Password = "sample-password"

	db, err := OpenPostgres(t.Context(), cfg)
	if db != nil {
		defer db.Close()
		t.Fatal("expected no pool on invalid config")
	}
	if err == nil || strings.Contains(err.Error(), cfg.Password) {
		t.Fatal("expected an error without the database password")
	}
}

func TestOpenPostgresCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	db, err := OpenPostgres(ctx, config.NewConfig().Postgres)
	if db != nil {
		defer db.Close()
		t.Fatal("expected no pool when the context is canceled")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
