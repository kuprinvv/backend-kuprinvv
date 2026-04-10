package room

import (
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/room/converter"
	"test-backend-1-kuprinvv/pkg/httpx"
)

// ListRooms godoc
//
//	@Summary     Список переговорок
//	@Description Возвращает список всех переговорок
//	@Tags        Переговорки
//	@Produce     json
//	@Success     200 {object} dto.ListRoomsResponse
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Security    BearerAuth
//	@Router      /rooms/list [get]
func (h *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomService.ListRooms(r.Context())
	if err != nil {
		httpx.Error(w, handler.ErrInternal, "failed to list rooms", http.StatusInternalServerError)
		return
	}

	httpx.JSON(w, converter.ServiceToListRoomsResponse(rooms), http.StatusOK)
}
