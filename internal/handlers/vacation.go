package handlers

import (
	"context"
	"errors"
	"sugdio/api"
	"sugdio/internal/domain"
)

func (h *Handler) CreateVacation(ctx context.Context, request api.CreateVacationRequestObject) (api.CreateVacationResponseObject, error) {
	vacation := domain.VacationCreate{
		EmployeeID: *request.Body.EmployeeId,
		StartDate:  request.Body.StartDate.Time,
		EndDate:    request.Body.EndDate.Time,
	}
	if request.Body.Comment != nil {
		vacation.Comment = request.Body.Comment
	}

	v, err := h.vacService.CreateVacation(ctx, &vacation)
	if err != nil {
		if errors.Is(err, domain.ErrDateOverlap) {
			return api.CreateVacation409JSONResponse{ConflictJSONResponse: api.ConflictJSONResponse{Code: "DATE_OVERLAP", Message: err.Error()}}, nil
		}
		return nil, err
	}
	return api.CreateVacation201JSONResponse(toApiVacation(v)), nil
}

func (h *Handler) DeleteVacation(ctx context.Context, request api.DeleteVacationRequestObject) (api.DeleteVacationResponseObject, error) {
	err := h.vacService.DeleteVacation(ctx, request.VacationID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.DeleteVacation404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Vacation not found"}}, nil
		}
	}
	return api.DeleteVacation204Response(api.DeleteVacation204Response{}), nil
}

func (h *Handler) GetVacation(ctx context.Context, request api.GetVacationRequestObject) (api.GetVacationResponseObject, error) {
	v, err := h.vacService.GetByIDVacation(ctx, request.VacationID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.GetVacation404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Vacation not found"}}, nil
		}
		return nil, err
	}
	return api.GetVacation200JSONResponse(toApiVacation(v)), nil
}

func (h *Handler) GetVacations(ctx context.Context, request api.GetVacationsRequestObject) (api.GetVacationsResponseObject, error) {
	params := request.Params
	var filter domain.VacationFilter
	filter.Limit = 50
	filter.Offset = 0

	if params.Filter != nil {
		filter.EmployeeID = params.Filter.EmployeeId
		if params.Filter.Status != nil {
			status := domain.VacationStatus(*params.Filter.Status)
			filter.Status = &status
		}
		if params.Filter.DateFrom != nil {
			filter.FromDate = &params.Filter.DateFrom.Time
		}
		if params.Filter.DateTo != nil {
			filter.ToDate = &params.Filter.DateTo.Time
		}

		if params.Filter.Offset != nil {
			filter.Offset = *params.Filter.Offset
		}
		if params.Filter.Limit != nil {
			filter.Limit = *params.Filter.Limit
		}
	}

	vacList, err := h.vacService.ListVacation(ctx, filter)
	if err != nil {
		return nil, err
	}

	total := len(vacList)
	res := api.VacationList{Total: &total}
	items := make([]api.Vacation, 0, filter.Limit)
	for _, v := range vacList {
		items = append(items, toApiVacation(&v))
	}
	res.Items = &items

	return api.GetVacations200JSONResponse(res), nil
}

func (h *Handler) PatchVacation(ctx context.Context, request api.PatchVacationRequestObject) (api.PatchVacationResponseObject, error) {
	patch := domain.VacationPatch{
		Comment: request.Body.Comment,
		Status:  (*domain.VacationStatus)(request.Body.Status),
	}

	if request.Body.StartDate != nil {
		patch.StartDate = &request.Body.StartDate.Time
	}
	if request.Body.EndDate != nil {
		patch.EndDate = &request.Body.EndDate.Time
	}

	vac, err := h.vacService.UpdateVacation(ctx, request.VacationID, patch)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.PatchVacation404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{Code: "NOT_FOUND", Message: "Vacation not found"}}, nil
		}
		if errors.Is(err, domain.ErrDateOverlap) {
			return api.PatchVacation409JSONResponse{ConflictJSONResponse: api.ConflictJSONResponse{Code: "DATE_OVERLAP", Message: err.Error()}}, nil
		}
		return nil, err
	}
	return api.PatchVacation200JSONResponse(toApiVacation(vac)), nil
}
