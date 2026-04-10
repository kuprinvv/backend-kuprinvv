package booking

import (
	"context"
	"log"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/internal/service"
	"time"

	"github.com/google/uuid"
)

func (s *serv) CreateBooking(ctx context.Context, userID uuid.UUID, slotID uuid.UUID, createConferenceLink bool) (*model.Booking, error) {
	slot, err := s.slotRepo.GetSlotByID(ctx, slotID)
	if err != nil {
		return nil, err
	}

	if slot.StartTime.Before(time.Now().UTC()) {
		return nil, model.ErrPastSlot
	}

	booking := &model.Booking{
		ID:     uuid.New(),
		UserID: userID,
		SlotID: slot.ID,
	}

	if createConferenceLink {
		link, err := s.conferenceClient.CreateLink(ctx)
		if err != nil {
			log.Printf("conference service unavailable, booking will be created without link: %v", err)
		} else {
			booking.ConferenceLink = &link
		}
	}

	booking, err = s.bookingRepo.CreateBooking(ctx, *booking)
	if err != nil {
		if service.IsUniqueViolation(err) {
			return nil, model.ErrSlotAlreadyBooked
		}
		return nil, err
	}

	return booking, nil
}
