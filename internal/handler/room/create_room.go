package room

import (
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/room/converter"
	"test-backend-1-kuprinvv/internal/handler/room/dto"
	"test-backend-1-kuprinvv/pkg/httpx"
)

// CreateRoom godoc
//
//	@Summary     Создать переговорку
//	@Description Создаёт новую переговорку. Только для роли admin.
//	@Tags        Переговорки
//	@Accept      json
//	@Produce     json
//	@Param       body body dto.CreateRoomRequest true "Данные переговорки"
//	@Success     201 {object} dto.CreateRoomResponse
//	@Failure     400 {object} httpx.ErrorResponse "invalid request"
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Security    BearerAuth
//	@Router      /rooms/create [post]
func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.HandleBody[dto.CreateRoomRequest](r)
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid body", http.StatusBadRequest)
		return
	}

	room, err := h.roomService.CreateRoom(r.Context(), converter.CreateRoomRequestToService(*body))
	if err != nil {
		httpx.Error(w, handler.ErrInternal, "failed to create room", http.StatusInternalServerError)
		return
	}

	httpx.JSON(w, converter.ServiceToCreateRoomResponse(*room), http.StatusCreated)
}
