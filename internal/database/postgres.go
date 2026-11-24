package database

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"

	"brewjenkins/go-ai-insight-advisor/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// OpenPostgres opens and verifies a connection pool. The caller must close it.
func OpenPostgres(ctx context.Context, cfg config.PostgresConfig) (*sqlx.DB, error) {
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     net.JoinHostPort(cfg.Host, cfg.Port),
		Path:     "/" + cfg.DBName,
		RawPath:  "/" + url.PathEscape(cfg.DBName),
		RawQuery: url.Values{"sslmode": {cfg.SSLMode}}.Encode(),
	}
	connConfig, err := pgx.ParseConfig(dsn.String())
	if err != nil {
		return nil, errors.New("invalid PostgreSQL connection settings")
	}

	db := sqlx.NewDb(stdlib.OpenDB(*connConfig), "pgx")
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}
	return db, nil
}
