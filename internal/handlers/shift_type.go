package handlers

import (
	"context"
	"errors"
	"sugdio/api"
	"sugdio/internal/domain"
)

func (h *Handler) CreateShiftType(ctx context.Context, request api.CreateShiftTypeRequestObject) (api.CreateShiftTypeResponseObject, error) {
	st := domain.ShiftTypeCreate{Name: request.Body.Name}

	r, err := h.shiftService.CreateShiftType(ctx, &st)
	if err != nil {
		return nil, err
	}
	return api.CreateShiftType201JSONResponse(toApiShiftType(r)), err
}

func (h *Handler) DeleteShiftType(ctx context.Context, request api.DeleteShiftTypeRequestObject) (api.DeleteShiftTypeResponseObject, error) {
	err := h.shiftService.DeleteShiftType(ctx, request.ShiftTypeID)
	if errors.Is(err, domain.ErrShiftTypeInUse) {
		return api.DeleteShiftType409JSONResponse{ConflictJSONResponse: api.ConflictJSONResponse{Code: "CONFLICT", Message: err.Error()}}, nil
	}
	if errors.Is(err, domain.ErrNotFound) {
		return api.DeleteShiftType404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Shift type not found"}}, nil
	}
	return api.DeleteShiftType204Response{}, err
}

func (h *Handler) GetShiftType(ctx context.Context, request api.GetShiftTypeRequestObject) (api.GetShiftTypeResponseObject, error) {
	st, err := h.shiftService.GetByIDShiftType(ctx, request.ShiftTypeID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.GetShiftType404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Shift type not found"}}, nil
		}
	}
	return api.GetShiftType200JSONResponse(toApiShiftType(st)), nil
}

func (h *Handler) GetShiftTypes(ctx context.Context, request api.GetShiftTypesRequestObject) (api.GetShiftTypesResponseObject, error) {
	stList, err := h.shiftService.ListShiftType(ctx)
	if err != nil {
		return nil, err
	}

	total := len(stList)
	items := make([]api.ShiftType, 0, total)
	res := api.ShiftTypeList{Total: &total}
	for _, v := range stList {
		items = append(items, toApiShiftType(&v))
	}
	res.Items = &items

	return api.GetShiftTypes200JSONResponse(res), nil
}

func (h *Handler) PatchShiftType(ctx context.Context, request api.PatchShiftTypeRequestObject) (api.PatchShiftTypeResponseObject, error) {
	patch := domain.ShiftTypePatch{
		Name: request.Body.Name,
	}

	st, err := h.shiftService.UpdateShiftType(ctx, request.ShiftTypeID, patch)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.PatchShiftType404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Shift type not found"}}, nil
		}
		return nil, err
	}
	return api.PatchShiftType200JSONResponse(toApiShiftType(st)), nil
}
