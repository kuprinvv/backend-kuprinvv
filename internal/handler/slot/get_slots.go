package slot

import (
	"errors"
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/slot/converter"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/pkg/httpx"
	"time"
)

const (
	roomIDURLParam = "roomId"
	dateQueryParam = "date"
)

// GetSlots godoc
//
//	@Summary     Доступные слоты переговорки
//	@Description Возвращает список свободных слотов переговорки на указанную дату
//	@Tags        Слоты
//	@Produce     json
//	@Param       roomId path  string true "UUID переговорки"
//	@Param       date   query string true "Дата в формате YYYY-MM-DD"
//	@Success     200 {object} dto.GetSlotsResponse
//	@Failure     400 {object} httpx.ErrorResponse "invalid request"
//	@Failure     404 {object} httpx.ErrorResponse "room not found"
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Security    BearerAuth
//	@Router      /rooms/{roomId}/slots/list [get]
func (h *Handler) GetSlots(w http.ResponseWriter, r *http.Request) {
	roomID, err := httpx.ParseUUIDParam(r, roomIDURLParam)
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid roomId", http.StatusBadRequest)
		return
	}

	date, err := httpx.QueryParam(r, dateQueryParam, func(s string) (time.Time, error) {
		return time.Parse(time.DateOnly, s)
	})
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid date", http.StatusBadRequest)
		return
	}

	slots, err := h.slotServ.GetSlots(r.Context(), roomID, date)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRoomNotFound):
			httpx.Error(w, handler.ErrRoomNotFound, "room not found", http.StatusNotFound)
		default:
			httpx.Error(w, handler.ErrInternal, "failed get slots: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	httpx.JSON(w, converter.ServiceToGetSlotsResponse(slots), http.StatusOK)
}
