package handlers

import (
	"context"
	"errors"
	"sugdio/api"
	"sugdio/internal/domain"
)

func (h *Handler) CreateShift(ctx context.Context, request api.CreateShiftRequestObject) (api.CreateShiftResponseObject, error) {
	shift := domain.ShiftCreate{
		Date:   request.Body.Date.Time,
		Status: (*domain.ShiftStatus)(request.Body.Status),
	}
	if request.Body.EmployeeId == nil {
		shift.EmployeeID = 0
	} else {
		shift.EmployeeID = int(*request.Body.EmployeeId)
	}
	if request.Body.ShiftTypeId == nil {
		shift.ShiftTypeID = 0
	} else {
		shift.ShiftTypeID = *request.Body.ShiftTypeId
	}

	s, err := h.shiftService.CreateShift(ctx, &shift)
	if err != nil {
		if errors.Is(err, domain.ErrWrongShiftStatus) || errors.Is(err, domain.ErrWrongShiftType) || errors.Is(err, domain.ErrWrongEmployee) {
			return api.CreateShift400JSONResponse{BadRequestJSONResponse: api.BadRequestJSONResponse{Code: "BAD_REQUEST", Message: err.Error()}}, nil
		}
		return nil, err
	}

	return api.CreateShift201JSONResponse(toAPIShift(s)), nil
}

func (h *Handler) DeleteShift(ctx context.Context, request api.DeleteShiftRequestObject) (api.DeleteShiftResponseObject, error) {
	err := h.shiftService.DeleteShift(ctx, request.ShiftID)
	if errors.Is(err, domain.ErrNotFound) {
		return api.DeleteShift404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Shift not found"}}, nil
	}
	return api.DeleteShift204Response{}, err
}

func (h *Handler) GetShift(ctx context.Context, request api.GetShiftRequestObject) (api.GetShiftResponseObject, error) {
	shift, err := h.shiftService.GetByIDShift(ctx, request.ShiftID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.GetShift404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Shift not found"}}, nil
		}
		return nil, err
	}
	return api.GetShift200JSONResponse(toAPIShift(shift)), nil
}

func (h *Handler) GetShifts(ctx context.Context, request api.GetShiftsRequestObject) (api.GetShiftsResponseObject, error) {
	params := request.Params
	var filter domain.ShiftFilter
	filter.Limit = 50
	filter.Offset = 0

	if params.Filter != nil {
		filter.DateFrom = &params.Filter.DateFrom.Time
		filter.DateTo = &params.Filter.DateTo.Time
		filter.EmployeeID = params.Filter.EmployeeId
		filter.ShiftTypeID = params.Filter.ShiftTypeId
		filter.Status = (*domain.ShiftStatus)(params.Filter.Status)

		if params.Filter.Offset != nil {
			filter.Offset = *params.Filter.Offset
		}
		if params.Filter.Limit != nil {
			filter.Limit = *params.Filter.Limit
		}
	}

	shiftList, err := h.shiftService.ListShift(ctx, filter)
	if err != nil {
		return nil, err
	}

	total := len(shiftList)
	res := api.ShiftList{Total: &total}
	items := make([]api.Shift, 0, filter.Limit)
	for _, v := range shiftList {
		items = append(items, toAPIShift(&v))
	}
	res.Items = &items
	return api.GetShifts200JSONResponse(res), nil
}

func (h *Handler) PatchShift(ctx context.Context, request api.PatchShiftRequestObject) (api.PatchShiftResponseObject, error) {
	patch := domain.ShiftPatch{
		ShiftTypeID: request.Body.ShiftTypeId,
		Date:        &request.Body.Date.Time,
		Status:      (*domain.ShiftStatus)(request.Body.Status),
		EmployeeID:  request.Body.EmployeeId,
	}

	shift, err := h.shiftService.UpdateShift(ctx, request.ShiftID, patch)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.PatchShift404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Shift not found"}}, nil
		}
		if errors.Is(err, domain.ErrWrongShiftStatus) || errors.Is(err, domain.ErrWrongShiftType) || errors.Is(err, domain.ErrWrongEmployee) {
			return api.PatchShift400JSONResponse{BadRequestJSONResponse: api.BadRequestJSONResponse{Code: "BAD_REQUEST", Message: err.Error()}}, nil
		}
		return nil, err
	}
	return api.PatchShift200JSONResponse(toAPIShift(shift)), nil
}
