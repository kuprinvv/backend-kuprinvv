package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	SlotID               uuid.UUID `json:"slotId" validate:"required"`
	CreateConferenceLink bool      `json:"createConferenceLink"`
}

type BookingResponse struct {
	Booking Booking `json:"booking"`
}

type BookingsResponse struct {
	Bookings []Booking `json:"bookings"`
}

type BookingsListResponse struct {
	Bookings   []Booking  `json:"bookings"`
	Pagination Pagination `json:"pagination"`
}

type Booking struct {
	ID             uuid.UUID  `json:"id"`
	SlotID         uuid.UUID  `json:"slotId"`
	UserID         uuid.UUID  `json:"userId"`
	Status         string     `json:"status"`
	ConferenceLink *string    `json:"conferenceLink,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
