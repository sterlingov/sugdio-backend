package service

import (
	"context"
	"sugdio/internal/domain"
	"sugdio/internal/repository"
)

type EmployeeService struct {
	repository           repository.EmployeeRepository
	departmentRepository repository.DepartmentRepository
	positionRepository   repository.PositionRepository
	userRepository       repository.UserRepository
	shiftRepository      repository.ShiftRepository
	vacationRepository   repository.VacationRepository
}

func NewEmployeeService(repo repository.EmployeeRepository, dr repository.DepartmentRepository, pr repository.PositionRepository, ur repository.UserRepository, sr repository.ShiftRepository, vr repository.VacationRepository) *EmployeeService {
	return &EmployeeService{repository: repo, departmentRepository: dr, positionRepository: pr, userRepository: ur, shiftRepository: sr, vacationRepository: vr}
}

func (e *EmployeeService) Create(ctx context.Context, emp *domain.EmployeeCreate) (*domain.Employee, error) {
	var employee *domain.Employee

	if emp.DepartmentId != nil {
		if _, err := e.departmentRepository.GetByIDDepartment(ctx, *emp.DepartmentId); err != nil {
			return employee, domain.ErrWrongDepartment
		}
	}

	if emp.PositionId != nil {
		if _, err := e.positionRepository.GetByIDPosition(ctx, *emp.PositionId); err != nil {
			return employee, domain.ErrWrongPosition
		}
	}

	if emp.UserId != nil {
		if _, err := e.userRepository.GetByIDUser(ctx, *emp.UserId); err != nil {
			return employee, domain.ErrWrongUser
		}
	}

	id, err := e.repository.CreateEmployee(ctx, emp)
	if err != nil {
		return employee, err
	}
	return e.repository.GetByIDEmployee(ctx, id)
}

func (e *EmployeeService) Delete(ctx context.Context, id int64) error {
	if hasShifts, _ := e.shiftRepository.HasShifts(ctx, id); hasShifts {
		return domain.ErrEmployeeHasShifts
	}

	if hasVacations, _ := e.vacationRepository.HasVacations(ctx, id); hasVacations {
		return domain.ErrEmployeeHasVacations
	}

	err := e.repository.DeleteEmployee(ctx, id)
	return err
}

func (e *EmployeeService) GetByID(ctx context.Context, id int64) (*domain.Employee, error) {
	return e.repository.GetByIDEmployee(ctx, id)
}

func (e *EmployeeService) List(ctx context.Context, filter domain.EmployeeFilter) ([]domain.Employee, error) {
	return e.repository.ListEmployee(ctx, filter)
}

func (e *EmployeeService) Update(ctx context.Context, id int64, patch domain.EmployeePatch) (*domain.Employee, error) {
	var emp *domain.Employee

	if patch.DepartmentId != nil {
		if _, err := e.departmentRepository.GetByIDDepartment(ctx, *patch.DepartmentId); err != nil {
			return emp, domain.ErrWrongDepartment
		}
	}

	if patch.PositionId != nil {
		if _, err := e.positionRepository.GetByIDPosition(ctx, *patch.PositionId); err != nil {
			return emp, domain.ErrWrongPosition
		}
	}

	if patch.UserId != nil {
		if _, err := e.userRepository.GetByIDUser(ctx, *patch.UserId); err != nil {
			return emp, domain.ErrWrongUser
		}
	}

	err := e.repository.UpdateEmployee(ctx, id, patch)
	if err != nil {
		return emp, err
	}
	return e.repository.GetByIDEmployee(ctx, id)
}
