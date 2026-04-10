package booking

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/google/uuid"
)

func (s *serv) CancelBooking(ctx context.Context, userID uuid.UUID, bookingID uuid.UUID) (*model.Booking, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.UserID != userID {
		return nil, model.ErrForbidden
	}

	if booking.Status == "cancelled" {
		return booking, nil
	}

	return s.bookingRepo.CancelBooking(ctx, bookingID)
}
