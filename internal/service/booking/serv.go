package booking

import (
	"test-backend-1-kuprinvv/internal/client"
	"test-backend-1-kuprinvv/internal/repository"
	"test-backend-1-kuprinvv/internal/service"
)

var _ service.BookingService = (*serv)(nil)

type serv struct {
	bookingRepo      repository.BookingRepository
	slotRepo         repository.SlotRepository
	conferenceClient client.ConferenceClient
}

func NewBookingService(bookingRepo repository.BookingRepository, slotRepo repository.SlotRepository, conferenceClient client.ConferenceClient) *serv {
	return &serv{
		bookingRepo:      bookingRepo,
		slotRepo:         slotRepo,
		conferenceClient: conferenceClient,
	}
}
