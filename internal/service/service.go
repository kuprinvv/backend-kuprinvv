package service

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	DummyLogin(ctx context.Context, userUUID uuid.UUID, role string) (string, error)
	Register(ctx context.Context, email, password, role string) (model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type SlotService interface {
	GetSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]model.Slot, error)
}

type BookingService interface {
	CreateBooking(ctx context.Context, userID uuid.UUID, slotID uuid.UUID, createConferenceLink bool) (*model.Booking, error)
	CancelBooking(ctx context.Context, userID uuid.UUID, bookingID uuid.UUID) (*model.Booking, error)
	GetMyBookings(ctx context.Context, userID uuid.UUID) ([]model.Booking, error)
	ListBookings(ctx context.Context, page, pageSize int) ([]model.Booking, *model.Pagination, error)
}

type RoomService interface {
	CreateRoom(ctx context.Context, room model.Room) (*model.Room, error)
	ListRooms(ctx context.Context) ([]model.Room, error)
}

type ScheduleService interface {
	CreateSchedule(ctx context.Context, schedule model.Schedule) (*model.Schedule, error)
}
