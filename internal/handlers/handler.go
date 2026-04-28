package handlers

import (
	"sugdio/internal/service"
)

type Handler struct {
	empService   *service.EmployeeService
	shiftService *service.ShiftService
}

func NewHandler(es *service.EmployeeService, ss *service.ShiftService) *Handler {
	return &Handler{empService: es, shiftService: ss}
}
