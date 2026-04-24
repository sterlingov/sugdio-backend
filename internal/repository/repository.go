package repository

import (
	"context"
	"sugdio/internal/domain"
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, emp *domain.EmployeeCreate) (int64, error)

	GetByIDEmployee(ctx context.Context, id int64) (*domain.Employee, error)

	ListEmployee(ctx context.Context, filter domain.EmployeeFilter) ([]domain.Employee, error)

	UpdateEmployee(ctx context.Context, id int64, patch domain.EmployeePatch) error

	DeleteEmployee(ctx context.Context, id int64) error
}

type AuthRepository interface {
	GetByEmailAuth(ctx context.Context, email string) (domain.UserCredentials, error)
}

type DepartmentRepository interface {
	GetByIDDepartment(ctx context.Context, id int) (*domain.Department, error)
	CreateDepartment(ctx context.Context, department *domain.DepartmentCreate) (int, error)
}

type PositionRepository interface {
	GetByIDPosition(ctx context.Context, id int) (*domain.Position, error)
	CreatePosition(ctx context.Context, position *domain.PositionCreate) (int, error)
}

type UserRepository interface {
	GetByIDUser(ctx context.Context, id int64) (*domain.UserShort, error)
}

type ShiftRepository interface {
	HasShifts(ctx context.Context, EmployeeID int64) (bool, error)
}

type VacationRepository interface {
	HasVacations(ctx context.Context, EmployeeID int64) (bool, error)
}
