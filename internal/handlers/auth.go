package handlers

import (
	"context"
	"errors"
	"sugdio/api"
	"sugdio/internal/domain"
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
