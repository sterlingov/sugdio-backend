package service

import (
	"context"
	"sugdio/internal/domain"
	"sugdio/internal/repository"
)

type ShiftService struct {
	repository         repository.ShiftRepository
	employeeRepository repository.EmployeeRepository
}

func NewShiftService(repo repository.ShiftRepository, er repository.EmployeeRepository) *ShiftService {
	return &ShiftService{repository: repo, employeeRepository: er}
}

func (s *ShiftService) CreateShift(ctx context.Context, sc *domain.ShiftCreate) (*domain.Shift, error) {
	var shift *domain.Shift

	if sc.Status != nil {
		if *sc.Status != domain.ShiftStatusCancelled && *sc.Status != domain.ShiftStatusCompleted && *sc.Status != domain.ShiftStatusPlanned {
			return nil, domain.ErrWrongShiftStatus
		}
	}

	if _, err := s.repository.GetByIDShiftType(ctx, sc.ShiftTypeID); err != nil {
		return nil, domain.ErrWrongShiftType
	}

	if _, err := s.employeeRepository.GetByIDEmployee(ctx, int64(sc.EmployeeID)); err != nil {
		return nil, domain.ErrWrongEmployee
	}

	id, err := s.repository.CreateShift(ctx, sc)
	if err != nil {
		return shift, err
	}
	return s.repository.GetByIDShift(ctx, id)
}

func (s *ShiftService) DeleteShift(ctx context.Context, id int) error {
	err := s.repository.DeleteShift(ctx, id)
	return err
}

func (s *ShiftService) GetByIDShift(ctx context.Context, id int) (*domain.Shift, error) {
	return s.repository.GetByIDShift(ctx, id)
}

func (s *ShiftService) ListShift(ctx context.Context, filter domain.ShiftFilter) ([]domain.Shift, error) {
	return s.repository.ListShift(ctx, filter)
}

func (s *ShiftService) UpdateShift(ctx context.Context, id int, patch domain.ShiftPatch) (*domain.Shift, error) {
	var shift *domain.Shift

	if patch.EmployeeID != nil {
		if _, err := s.employeeRepository.GetByIDEmployee(ctx, int64(*patch.EmployeeID)); err != nil {
			return nil, domain.ErrWrongEmployee
		}
	}
	if patch.ShiftTypeID != nil {
		if _, err := s.repository.GetByIDShiftType(ctx, *patch.ShiftTypeID); err != nil {
			return nil, domain.ErrWrongShiftType
		}
	}
	if patch.Status != nil {
		if *patch.Status != domain.ShiftStatusCancelled && *patch.Status != domain.ShiftStatusCompleted && *patch.Status != domain.ShiftStatusPlanned {
			return nil, domain.ErrWrongShiftStatus
		}
	}

	err := s.repository.UpdateShift(ctx, id, patch)
	if err != nil {
		return shift, err
	}
	return s.repository.GetByIDShift(ctx, id)
}

func (s *ShiftService) CreateShiftType(ctx context.Context, st *domain.ShiftTypeCreate) (*domain.ShiftType, error) {
	return s.CreateShiftType(ctx, st)
}
func (s *ShiftService) DeleteShiftType(ctx context.Context, id int) error {
	filter := domain.ShiftFilter{ShiftTypeID: &id, Limit: 1, Offset: 0}
	shifts, err := s.repository.ListShift(ctx, filter)
	if len(shifts) == 0 {
		return domain.ErrShiftTypeInUse
	}
	return err
}
func (s *ShiftService) UpdateShiftType(ctx context.Context, id int, patch domain.ShiftTypePatch) (*domain.ShiftType, error) {
	err := s.repository.UpdateShiftType(ctx, id, patch)
	if err != nil {
		return nil, err
	}
	return s.repository.GetByIDShiftType(ctx, id)
}
func (s *ShiftService) ListShiftType(ctx context.Context) ([]domain.ShiftType, error) {
	return s.repository.ListShiftType(ctx)
}
func (s *ShiftService) GetByIDShiftType(ctx context.Context, id int) (*domain.ShiftType, error) {
	return s.repository.GetByIDShiftType(ctx, id)
}
