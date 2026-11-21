package employee

import "errors"

var (
	ErrNotFound   = errors.New("employee not found")
	ErrEmailTaken = errors.New("email already exists")
)

type Employee struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
