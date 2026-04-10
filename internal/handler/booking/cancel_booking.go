package booking

import (
	"errors"
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/booking/converter"
	"test-backend-1-kuprinvv/internal/middleware"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/pkg/httpx"
)

const (
	bookingIdURLParam = "bookingId"
)

// CancelBooking godoc
//
//	@Summary     Отменить бронь
//	@Description Отменяет бронь текущего пользователя. Операция идемпотентна. Только для роли user.
//	@Tags        Брони
//	@Produce     json
//	@Param       bookingId path string true "UUID брони"
//	@Success     200 {object} dto.BookingResponse
//	@Failure     400 {object} httpx.ErrorResponse "invalid UUID"
//	@Failure     403 {object} httpx.ErrorResponse "cannot cancel another user's booking"
//	@Failure     404 {object} httpx.ErrorResponse "booking not found"
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Security    BearerAuth
//	@Router      /bookings/{bookingId}/cancel [post]
func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	bookingID, err := httpx.ParseUUIDParam(r, bookingIdURLParam)
	if err != nil {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid booking id", http.StatusBadRequest)
		return
	}

	booking, err := h.bookingServ.CancelBooking(r.Context(), user.UserID, bookingID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrForbidden):
			httpx.Error(w, handler.ErrForbidden,
				"cannot cancel another user's booking", http.StatusForbidden)
		case errors.Is(err, model.ErrBookingNotFound):
			httpx.Error(w, handler.ErrBookingNotFound,
				"booking not found", http.StatusNotFound)
		default:
			httpx.Error(w, handler.ErrInternal, "failed cancelled booking", http.StatusInternalServerError)
		}
		return
	}

	httpx.JSON(w, converter.ServiceToBookingResponse(*booking), http.StatusOK)
}
