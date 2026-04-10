package booking

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/google/uuid"
)

func (s *serv) GetMyBookings(ctx context.Context, userID uuid.UUID) ([]model.Booking, error) {
	return s.bookingRepo.GetFutureBookingsByUser(ctx, userID)
}
