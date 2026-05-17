package handlers

import (
	"sugdio/internal/service"
)

type Handler struct {
	authService  *service.AuthService
	empService   *service.EmployeeService
	shiftService *service.ShiftService
	vacService   *service.VacationService
}

func NewHandler(as *service.AuthService, es *service.EmployeeService, ss *service.ShiftService, vs *service.VacationService) *Handler {
	return &Handler{empService: es, shiftService: ss, vacService: vs, authService: as}
}
