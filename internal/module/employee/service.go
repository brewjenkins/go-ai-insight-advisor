package employee

import (
	"context"
	"fmt"
	"log/slog"
)

// EmployeeRepository is the consumer-side contract required by Service.
type EmployeeRepository interface {
	Create(ctx context.Context, emp *Employee) error
	FindByID(ctx context.Context, id int64) (*Employee, error)
}

type Service struct {
	repo   EmployeeRepository
	logger *slog.Logger
}

func NewService(repo EmployeeRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*Employee, error) {
	emp := &Employee{
		Name:  req.Name,
		Email: req.Email,
	}
	if err := s.repo.Create(ctx, emp); err != nil {
		s.logger.Error("failed to create employee", "error", err)
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
