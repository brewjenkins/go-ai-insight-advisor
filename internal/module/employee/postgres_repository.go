package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, emp *Employee) error {
	const query = `INSERT INTO employees (name, email) VALUES ($1, $2) RETURNING id`

	if err := r.db.QueryRowContext(ctx, query, emp.Name, emp.Email).Scan(&emp.ID); err != nil {
		return fmt.Errorf("insert employee: %w", err)
	}
	return nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id int64) (*Employee, error) {
	const query = `SELECT id, name, email FROM employees WHERE id = $1`

	var emp Employee
	err := r.db.QueryRowContext(ctx, query, id).Scan(&emp.ID, &emp.Name, &emp.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find employee by id: %w", err)
	}
	return &emp, nil
}
