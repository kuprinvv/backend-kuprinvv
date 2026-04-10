package auth

import (
	"errors"
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/auth/converter"
	"test-backend-1-kuprinvv/internal/handler/auth/dto"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/pkg/httpx"
)

// Register godoc
//
//	@Summary     Регистрация пользователя
//	@Description Создаёт нового пользователя по email и паролю, возвращает данные пользователя
//	@Tags        Авторизация
//	@Accept      json
//	@Produce     json
//	@Param       body body dto.RegisterRequest true "Данные для регистрации"
//	@Success     201 {object} dto.RegisterResponse
//	@Failure     400 {object} httpx.ErrorResponse "invalid request"
//	@Failure     400 {object} httpx.ErrorResponse "email already taken"
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Router      /register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.HandleBody[dto.RegisterRequest](r)
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), body.Email, body.Password, body.Role)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			httpx.Error(w, handler.ErrInvalidRequest, "user already exists", http.StatusBadRequest)
			return
		}
		httpx.Error(w, handler.ErrInternal, "internal error", http.StatusInternalServerError)
		return
	}

	httpx.JSON(w, converter.ServiceToRegisterResponse(user), http.StatusCreated)
}
