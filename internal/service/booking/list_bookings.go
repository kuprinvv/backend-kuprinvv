package booking

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"
)

func (s *serv) ListBookings(ctx context.Context, page, pageSize int) ([]model.Booking, *model.Pagination, error) {
	offset := (page - 1) * pageSize
	bookings, total, err := s.bookingRepo.ListBookings(ctx, pageSize, offset)
	if err != nil {
		return nil, nil, err
	}
	return bookings, &model.Pagination{Total: total, Page: page, PageSize: pageSize}, nil
}
