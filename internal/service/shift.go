package service

import "sugdio/internal/repository"

type ShiftService struct {
	repository repository.ShiftRepository
}

func NewShiftService(repo repository.ShiftRepository) *ShiftService {
	return &ShiftService{repository: repo}
}
