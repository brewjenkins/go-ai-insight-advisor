package employee

import (
	"errors"
	"time"
)

var (
	ErrNotFound   = errors.New("employee not found")
	ErrEmailTaken = errors.New("email already exists")
)

type Employee struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
