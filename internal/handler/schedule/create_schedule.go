package schedule

import (
	"errors"
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/schedule/converter"
	"test-backend-1-kuprinvv/internal/handler/schedule/dto"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/pkg/httpx"
)

const (
	roomIdParam = "roomId"
)

// CreateSchedule godoc
//
//	@Summary     Создать расписание переговорки
//	@Description Создаёт расписание доступности для переговорки. Расписание можно задать только один раз. Только для роли admin.
//	@Tags        Расписание
//	@Accept      json
//	@Produce     json
//	@Param       roomId path string                    true "UUID переговорки"
//	@Param       body   body dto.CreateScheduleRequest true "Расписание"
//	@Success     201 {object} dto.CreateScheduleResponse
//	@Failure     400 {object} httpx.ErrorResponse "invalid request"
//	@Failure     404 {object} httpx.ErrorResponse "room not found"
//	@Failure     409 {object} httpx.ErrorResponse "schedule already exists"
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Security    BearerAuth
//	@Router      /rooms/{roomId}/schedule/create [post]
func (h *Handler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	roomID, err := httpx.ParseUUIDParam(r, roomIdParam)
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid roomId", http.StatusBadRequest)
		return
	}

	body, err := httpx.HandleBody[dto.CreateScheduleRequest](r)
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid request", http.StatusBadRequest)
		return
	}

	body.RoomID = roomID

	svcSchedule, err := converter.CreateScheduleCreateToService(*body)
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, err.Error(), http.StatusBadRequest)
		return
	}

	schedule, err := h.scheduleServ.CreateSchedule(r.Context(), svcSchedule)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRoomNotFound):
			httpx.Error(w, handler.ErrRoomNotFound, "room not found", http.StatusNotFound)
		case errors.Is(err, model.ErrScheduleAlreadyExists):
			httpx.Error(w, handler.ErrScheduleExists,
				"schedule for this room already exists and cannot be changed", http.StatusConflict)
		default:
			httpx.Error(w, handler.ErrInternal, "internal error", http.StatusInternalServerError)
		}
		return
	}

	httpx.JSON(w, converter.ServiceToCreateScheduleResponse(*schedule), http.StatusCreated)
}
