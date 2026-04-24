package handlers

import (
	"context"
	"errors"
	"sugdio/api"
	"sugdio/internal/domain"
	"sugdio/internal/service"
)

type Handler struct {
	empService   *service.EmployeeService
	shiftService *service.ShiftService
}

func NewHandler(es *service.EmployeeService, ss *service.ShiftService) *Handler {
	return &Handler{empService: es, shiftService: ss}
}

func (h *Handler) CreateEmployee(ctx context.Context, request api.CreateEmployeeRequestObject) (api.CreateEmployeeResponseObject, error) {
	emp := domain.EmployeeCreate{
		FirstName:  request.Body.FirstName,
		MiddleName: request.Body.MiddleName,
		SecondName: request.Body.SecondName,
		Active:     request.Body.Active,

		DepartmentId: request.Body.DepartmentId,
		PositionId:   request.Body.PositionId,
		UserId:       request.Body.UserId,
	}

	employee, err := h.empService.Create(ctx, &emp)
	if err != nil {
		if errors.Is(err, domain.ErrWrongPosition) || errors.Is(err, domain.ErrWrongDepartment) || errors.Is(err, domain.ErrWrongUser) {
			return api.CreateEmployee400JSONResponse{BadRequestJSONResponse: api.BadRequestJSONResponse{Code: "BAD_REQUEST", Message: err.Error()}}, nil
		}
		return nil, err
	}

	return api.CreateEmployee201JSONResponse(toAPIEmployee(employee)), nil
}

func (h *Handler) DeleteEmployee(ctx context.Context, request api.DeleteEmployeeRequestObject) (api.DeleteEmployeeResponseObject, error) {
	err := h.empService.Delete(ctx, request.EmployeeID)
	if errors.Is(err, domain.ErrEmployeeHasShifts) || errors.Is(err, domain.ErrEmployeeHasVacations) {
		return api.DeleteEmployee409JSONResponse{ConflictJSONResponse: api.ConflictJSONResponse{Code: "CONFLICT", Message: err.Error()}}, nil
	}

	return api.DeleteEmployee204Response{}, err
}

func (h *Handler) GetEmployee(ctx context.Context, request api.GetEmployeeRequestObject) (api.GetEmployeeResponseObject, error) {
	emp, err := h.empService.GetByID(ctx, request.EmployeeID)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.GetEmployee404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Employee not found"}}, nil
		}
		return nil, err
	}

	return api.GetEmployee200JSONResponse(toAPIEmployee(emp)), nil
}

func (h *Handler) GetEmployees(ctx context.Context, request api.GetEmployeesRequestObject) (api.GetEmployeesResponseObject, error) {
	params := request.Params

	var filter domain.EmployeeFilter
	filter.Limit = 50
	filter.Offset = 0

	if params.Filter != nil {
		filter.FirstName = params.Filter.FirstName
		filter.SecondName = params.Filter.SecondName
		filter.Active = params.Filter.Active
		filter.DepartmentId = params.Filter.DepartmentId

		if params.Filter.Offset != nil {
			filter.Offset = *params.Filter.Offset
		}
		if params.Filter.Limit != nil {
			filter.Limit = *params.Filter.Limit
		}
	}

	empList, err := h.empService.List(ctx, filter)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.GetEmployees404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "No employees found with provided filter"}}, nil
		}
		return nil, err
	}

	total := len(empList)

	res := api.EmployeeList{Total: &total}
	items := make([]api.Employee, 0, total)
	for _, v := range empList {
		items = append(items, toAPIEmployee(&v))
	}
	res.Items = &items
	return api.GetEmployees200JSONResponse(res), err
}

func (h *Handler) PatchEmployee(ctx context.Context, request api.PatchEmployeeRequestObject) (api.PatchEmployeeResponseObject, error) {
	patch := domain.EmployeePatch{
		Active:       request.Body.Active,
		DepartmentId: request.Body.DepartmentId,
		FirstName:    request.Body.FirstName,
		MiddleName:   request.Body.MiddleName,
		PositionId:   request.Body.PositionId,
		SecondName:   request.Body.SecondName,
		UserId:       request.Body.UserId,
	}

	emp, err := h.empService.Update(ctx, request.EmployeeID, patch)
	if err != nil {
		if errors.Is(err, domain.ErrWrongPosition) || errors.Is(err, domain.ErrWrongDepartment) || errors.Is(err, domain.ErrWrongUser) {
			return api.PatchEmployee400JSONResponse{BadRequestJSONResponse: api.BadRequestJSONResponse{Code: "BAD_REQUEST", Message: err.Error()}}, nil
		}
		if errors.Is(err, domain.ErrNotFound) {
			return api.PatchEmployee404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Employee not found"}}, nil
		}
		return nil, err
	}
	return api.PatchEmployee200JSONResponse(toAPIEmployee(emp)), err
}
