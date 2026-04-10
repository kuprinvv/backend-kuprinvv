package auth

import (
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/auth/dto"
	"test-backend-1-kuprinvv/pkg/httpx"

	"github.com/google/uuid"
)

const (
	adminUUID = "11111111-1111-1111-1111-111111111111"
	userUUID  = "22222222-2222-2222-2222-222222222222"
)

// DummyLogin godoc
//
//	@Summary     Получить тестовый JWT
//	@Description Выдаёт JWT-токен по роли. UUID пользователя фиксирован для каждой роли: один UUID для admin, другой для user.
//	@Tags        Авторизация
//	@Accept      json
//	@Produce     json
//	@Param       body body dto.DummyLoginRequest true "Роль пользователя"
//	@Success     200 {object} dto.DummyLoginResponse
//	@Failure     400 {object} httpx.ErrorResponse "invalid request"
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Router      /dummyLogin [post]
func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.HandleBody[dto.DummyLoginRequest](r)
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid request", http.StatusBadRequest)
		return
	}

	var userID uuid.UUID

	switch body.Role {
	case "admin":
		userID = uuid.MustParse(adminUUID)
	case "user":
		userID = uuid.MustParse(userUUID)
	}

	token, err := h.authService.DummyLogin(r.Context(), userID, body.Role)
	if err != nil {
		httpx.Error(w, handler.ErrInternal, "internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.JSON(w, dto.DummyLoginResponse{Token: token}, http.StatusOK)
}
