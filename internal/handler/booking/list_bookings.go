package booking

import (
	"net/http"
	"strconv"
	"test-backend-1-kuprinvv/internal/handler"
	"test-backend-1-kuprinvv/internal/handler/booking/converter"
	"test-backend-1-kuprinvv/pkg/httpx"
)

const (
	pageQueryParam     = "page"
	pageSizeQueryParam = "pageSize"
)

// ListBookings godoc
//
//	@Summary     Список всех броней
//	@Description Возвращает постраничный список всех броней. Только для роли admin.
//	@Tags        Брони
//	@Produce     json
//	@Param       page     query int false "Номер страницы (по умолчанию 1)"
//	@Param       pageSize query int false "Размер страницы (по умолчанию 20, максимум 100)"
//	@Success     200 {object} dto.BookingsListResponse
//	@Failure     400 {object} httpx.ErrorResponse "invalid pagination parameters"
//	@Failure     500 {object} httpx.ErrorResponse "internal error"
//	@Security    BearerAuth
//	@Router      /bookings/list [get]
func (h *Handler) ListBookings(w http.ResponseWriter, r *http.Request) {
	page, err := httpx.QueryParam(r, pageQueryParam, strconv.Atoi, 1)
	if err != nil || page < 1 {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid page", http.StatusBadRequest)
		return
	}

	pageSize, err := httpx.QueryParam(r, pageSizeQueryParam, strconv.Atoi, 20)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		httpx.Error(w, handler.ErrInvalidRequest, "invalid pageSize", http.StatusBadRequest)
		return
	}

	bookings, pagination, err := h.bookingServ.ListBookings(r.Context(), page, pageSize)
	if err != nil {
		httpx.Error(w, handler.ErrInternal, "internal server error", http.StatusInternalServerError)
		return
	}

	httpx.JSON(w, converter.ServiceToBookingsListResponse(bookings, *pagination), http.StatusOK)
}
