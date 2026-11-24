package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, emp *Employee) error {
	const query = `
		INSERT INTO employees (name, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
		RETURNING id`

	err := r.db.GetContext(ctx, emp, query, emp.Name, emp.Email, emp.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrEmailTaken
	}
	if err != nil {
		return fmt.Errorf("insert employee: %w", err)
	}
	return nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id int64) (*Employee, error) {
	const query = `SELECT id, name, email FROM employees WHERE id = $1`

	var emp Employee
	err := r.db.GetContext(ctx, &emp, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find employee by id: %w", err)
	}
	return &emp, nil
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*Employee, error) {
	const query = `SELECT id, name, email, password_hash FROM employees WHERE email = $1`

	var emp Employee
	err := r.db.GetContext(ctx, &emp, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find employee by email: %w", err)
	}
	return &emp, nil
}
