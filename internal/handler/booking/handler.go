package booking

import "test-backend-1-kuprinvv/internal/service"

type Handler struct {
	bookingServ service.BookingService
}

func NewBookingHandler(bookingServ service.BookingService) *Handler {
	return &Handler{
		bookingServ: bookingServ,
	}
}
