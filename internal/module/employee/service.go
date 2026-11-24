package employee

import (
	"context"
	"errors"
	"fmt"

	"brewjenkins/go-ai-insight-advisor/internal/password"
)

// EmployeeRepository is the consumer-side contract required by Service.
type EmployeeRepository interface {
	Create(ctx context.Context, emp *Employee) error
	FindByID(ctx context.Context, id int64) (*Employee, error)
	FindByEmail(ctx context.Context, email string) (*Employee, error)
}

type Service struct {
	repo EmployeeRepository
}

func NewService(repo EmployeeRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*Employee, error) {
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailTaken
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("check existing employee: %w", err)
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash employee password: %w", err)
	}

	emp := &Employee{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
	}

	if err := s.repo.Create(ctx, emp); err != nil {
		return nil, fmt.Errorf("create employee: %w", err)
	}
	return emp, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Employee, error) {
	emp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get employee %d: %w", id, err)
	}
	return emp, nil
}
