package service

import (
	"context"
	"sugdio/internal/domain"
	"sugdio/internal/repository"
)

type VacationService struct {
	repository      repository.VacationRepository
	shiftRepository repository.ShiftRepository
}

func NewVacationService(repo repository.VacationRepository, shiftRepo repository.ShiftRepository) *VacationService {
	return &VacationService{repository: repo, shiftRepository: shiftRepo}
}

func (s *VacationService) CreateVacation(ctx context.Context, sc *domain.VacationCreate) (*domain.Vacation, error) {
	var vacation *domain.Vacation

	id, err := s.repository.CreateVacation(ctx, sc)

	if err != nil {
		return vacation, err
	}
	return s.repository.GetByIDVacation(ctx, id)
}

func (s *VacationService) DeleteVacation(ctx context.Context, id int) error {
	return s.repository.DeleteVacation(ctx, id)
}

func (s *VacationService) GetByIDVacation(ctx context.Context, id int) (*domain.Vacation, error) {
	return s.repository.GetByIDVacation(ctx, id)
}

func (s *VacationService) ListVacation(ctx context.Context, filter domain.VacationFilter) ([]domain.Vacation, error) {
	return s.repository.ListVacation(ctx, filter)
}

func (s *VacationService) UpdateVacation(ctx context.Context, id int, patch domain.VacationPatch) (*domain.Vacation, error) {
	var vacation *domain.Vacation

	err := s.repository.UpdateVacation(ctx, id, patch)
	if err != nil {
		return vacation, err
	}
	return s.repository.GetByIDVacation(ctx, id)
}
