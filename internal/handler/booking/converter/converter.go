package converter

import (
	"test-backend-1-kuprinvv/internal/handler/booking/dto"
	"test-backend-1-kuprinvv/internal/model"
)

func ServiceToBookingResponse(booking model.Booking) dto.BookingResponse {
	return dto.BookingResponse{
		Booking: serviceBookingToDtoBooking(booking),
	}
}

func ServiceToBookingsResponse(bookings []model.Booking) dto.BookingsResponse {
	resp := dto.BookingsResponse{Bookings: make([]dto.Booking, 0, len(bookings))}
	for _, booking := range bookings {
		resp.Bookings = append(resp.Bookings, serviceBookingToDtoBooking(booking))
	}
	return resp
}

func ServiceToBookingsListResponse(bookings []model.Booking, pagination model.Pagination) dto.BookingsListResponse {
	resp := dto.BookingsListResponse{Bookings: make([]dto.Booking, 0, len(bookings))}
	for _, booking := range bookings {
		resp.Bookings = append(resp.Bookings, serviceBookingToDtoBooking(booking))
	}
	resp.Pagination = dto.Pagination{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
		Total:    pagination.Total,
	}
	return resp
}

func serviceBookingToDtoBooking(booking model.Booking) dto.Booking {
	return dto.Booking{
		ID:             booking.ID,
		SlotID:         booking.SlotID,
		UserID:         booking.UserID,
		Status:         booking.Status,
		ConferenceLink: booking.ConferenceLink,
		CreatedAt:      booking.CreatedAt,
	}
}
