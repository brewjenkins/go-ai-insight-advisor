package authentication

import (
	"context"
	"errors"
	"fmt"

	"brewjenkins/go-ai-insight-advisor/internal/module/employee"
	"brewjenkins/go-ai-insight-advisor/internal/password"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type EmployeeLookup interface {
	FindByEmail(ctx context.Context, email string) (*employee.Employee, error)
}

type Service struct {
	employees EmployeeLookup
}

func NewService(employees EmployeeLookup) *Service {
	return &Service{employees: employees}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) error {
	emp, err := s.employees.FindByEmail(ctx, req.Username)
	if errors.Is(err, employee.ErrNotFound) {
		return ErrInvalidCredentials
	}
	if err != nil {
		return fmt.Errorf("look up login credentials: %w", err)
	}
	if !password.Verify(emp.PasswordHash, req.Password) {
		return ErrInvalidCredentials
	}
	return nil
}
