package booking

import (
	"net/http"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/booking/converter"
	"test-backend-1-kuprinvv/internal/middleware"
	"test-backend-1-kuprinvv/pkg/httpx"
)

// GetMyBookings godoc
//
//	@Summary     Мои брони
//	@Description Возвращает список будущих броней текущего пользователя. Только для роли user.
//	@Tags        Брони
//	@Produce     json
//	@Success     200 {object} dto.BookingsResponse
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Security    BearerAuth
//	@Router      /bookings/my [get]
func (h *Handler) GetMyBookings(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	bookings, err := h.bookingServ.GetMyBookings(r.Context(), user.UserID)
	if err != nil {
		httpx.Error(w, handler.ErrInternal, "failed to get booking", http.StatusInternalServerError)
		return
	}

	httpx.JSON(w, converter.ServiceToBookingsResponse(bookings), http.StatusOK)
}
