package app

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"brewjenkins/go-ai-insight-advisor/internal/config"
)

func TestNewRouterHealthz(t *testing.T) {
	cfg := config.NewConfig()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	r := NewRouter(cfg, nil, logger)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /healthz status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Header().Get("X-Request-ID") == "" {
		t.Fatal("GET /healthz missing X-Request-ID header")
	}
}
