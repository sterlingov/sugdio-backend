package handlers

import (
	"sugdio/internal/service"
)

type Handler struct {
	empService   *service.EmployeeService
	shiftService *service.ShiftService
	vacService   *service.VacationService
}

func NewHandler(es *service.EmployeeService, ss *service.ShiftService, vs *service.VacationService) *Handler {
	return &Handler{empService: es, shiftService: ss, vacService: vs}
}
