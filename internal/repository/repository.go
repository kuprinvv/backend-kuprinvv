package repository

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user model.User) error
}

type RoomRepository interface {
	GetRoomByID(ctx context.Context, id uuid.UUID) (*model.Room, error)
	CreateRoom(ctx context.Context, room model.Room) (*model.Room, error)
	ListRooms(ctx context.Context) ([]model.Room, error)
}

type ScheduleRepository interface {
	CreateSchedule(ctx context.Context, schedule model.Schedule) error
	GetByRoom(ctx context.Context, roomID uuid.UUID) (*model.Schedule, error)
}

type SlotRepository interface {
	GetSlotsByRoomAndDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]model.Slot, error)
	GetSlotByID(ctx context.Context, id uuid.UUID) (*model.Slot, error)
	CreateSlots(ctx context.Context, slots []model.Slot) error
}

type BookingRepository interface {
	GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error)
	GetFutureBookingsByUser(ctx context.Context, userID uuid.UUID) ([]model.Booking, error)
	ListBookings(ctx context.Context, limit, offset int) ([]model.Booking, int, error)
	CreateBooking(ctx context.Context, booking model.Booking) (*model.Booking, error)
	CancelBooking(ctx context.Context, bookingID uuid.UUID) (*model.Booking, error)
}
