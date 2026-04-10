package booking

import (
	"errors"
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/booking/converter"
	"test-backend-1-kuprinvv/internal/handler/booking/dto"
	"test-backend-1-kuprinvv/internal/middleware"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/pkg/httpx"
)

// CreateBooking godoc
//
//	@Summary     Создать бронь
//	@Description Создаёт бронь на слот от имени текущего пользователя. Только для роли user.
//	@Tags        Брони
//	@Accept      json
//	@Produce     json
//	@Param       body body dto.CreateBookingRequest true "Данные брони"
//	@Success     201 {object} dto.BookingResponse
//	@Failure     400 {object} httpx.ErrorResponse "invalid request or slot in the past"
//	@Failure     404 {object} httpx.ErrorResponse "slot not found"
//	@Failure     409 {object} httpx.ErrorResponse "slot already booked"
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Security    BearerAuth
//	@Router      /bookings/create [post]
func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	body, err := httpx.HandleBody[dto.CreateBookingRequest](r)
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	booking, err := h.bookingServ.CreateBooking(r.Context(), user.UserID, body.SlotID, body.CreateConferenceLink)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrSlotAlreadyBooked):
			httpx.Error(w, handler.ErrSlotBooked, "slot is already booked", http.StatusConflict)
		case errors.Is(err, model.ErrPastSlot):
			httpx.Error(w, handler.ErrInvalidRequest, "cannot book slot", http.StatusBadRequest)
		case errors.Is(err, model.ErrSlotNotFound):
			httpx.Error(w, handler.ErrSlotNotFound, "slot not found", http.StatusNotFound)
		default:
			httpx.Error(w, handler.ErrInternal, "internal error", http.StatusInternalServerError)
		}
		return
	}

	httpx.JSON(w, converter.ServiceToBookingResponse(*booking), http.StatusCreated)
}
