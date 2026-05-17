package handlers

import (
	"context"
	"errors"
	"sugdio/api"
	"sugdio/internal/domain"

	"github.com/oapi-codegen/runtime/types"
)

func (h *Handler) AuthLogin(ctx context.Context, request api.AuthLoginRequestObject) (api.AuthLoginResponseObject, error) {
	token, err := h.authService.Login(ctx, string(request.Body.Email), request.Body.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return api.AuthLogin401JSONResponse{UnauthorizedJSONResponse: api.UnauthorizedJSONResponse{Code: "WRONG_CREDENTIALS", Message: err.Error()}}, nil
		}
		return nil, err
	}
	return api.AuthLogin200JSONResponse{Token: &token}, nil
}

func (h *Handler) AuthRegister(ctx context.Context, request api.AuthRegisterRequestObject) (api.AuthRegisterResponseObject, error) {
	email := request.Body.Email
	password := request.Body.Password
	role := request.Body.Role

	if len(password) < 8 {
		return api.AuthRegister400JSONResponse{BadRequestJSONResponse: api.BadRequestJSONResponse{Code: "BAD_REQUEST", Message: "password must be at least 8 characters long"}}, nil
	}

	r, err := h.authService.Register(ctx, string(email), password, role)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return api.AuthRegister409JSONResponse{ConflictJSONResponse: api.ConflictJSONResponse{Code: "CONFLICT", Message: "user with this email already exists"}}, nil
		}
		if errors.Is(err, domain.ErrWrongUserRole) {
			return api.AuthRegister400JSONResponse{BadRequestJSONResponse: api.BadRequestJSONResponse{Code: "BAD_REQUEST", Message: "Wrong user role or you don't have permission to create user with this role"}}, nil
		}
		return nil, err
	}

	ID := int(r.ID)

	return api.AuthRegister201JSONResponse{Email: (*types.Email)(&r.Email), Id: &ID, Role: &r.Role}, nil
}
