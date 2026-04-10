package auth

import (
	"errors"
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/auth/dto"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/pkg/httpx"
)

// Login godoc
//
//	@Summary     Вход по email и паролю
//	@Description Проверяет учётные данные и возвращает JWT-токен
//	@Tags        Авторизация
//	@Accept      json
//	@Produce     json
//	@Param       body body dto.LoginRequest true "Учётные данные"
//	@Success     200 {object} dto.DummyLoginResponse
//	@Failure     400 {object} httpx.ErrorResponse "invalid request"
//	@Failure     401 {object} httpx.ErrorResponse "invalid credentials"
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Router      /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.HandleBody[dto.LoginRequest](r)
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			httpx.Error(w, handler.ErrUnauthorized, "invalid credentials", http.StatusUnauthorized)
			return
		}
		httpx.Error(w, handler.ErrInternal, "internal error", http.StatusInternalServerError)
		return
	}

	httpx.JSON(w, dto.DummyLoginResponse{Token: token}, http.StatusOK)
}
