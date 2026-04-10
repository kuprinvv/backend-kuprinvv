package di

import (
	"context"
	"test-backend-1-kuprinvv/internal/service"
	"test-backend-1-kuprinvv/internal/service/auth"
	"test-backend-1-kuprinvv/internal/service/booking"
	"test-backend-1-kuprinvv/internal/service/room"
	"test-backend-1-kuprinvv/internal/service/schedule"
	"test-backend-1-kuprinvv/internal/service/slot"
)

type serviceProvider struct {
	auth     service.AuthService
	booking  service.BookingService
	room     service.RoomService
	schedule service.ScheduleService
	slot     service.SlotService
}

func (c *Container) AuthService(ctx context.Context) service.AuthService {
	if c.services.auth == nil {
		c.services.auth = auth.NewAuthService(c.JWTConfig(), c.UserRepo(ctx))
	}
	return c.services.auth
}

func (c *Container) BookingService(ctx context.Context) service.BookingService {
	if c.services.booking == nil {
		c.services.booking = booking.NewBookingService(c.BookingRepo(ctx), c.SlotRepo(ctx), c.ConferenceClient())
	}
	return c.services.booking
}

func (c *Container) RoomService(ctx context.Context) service.RoomService {
	if c.services.room == nil {
		c.services.room = room.NewRoomService(c.RoomRepo(ctx))
	}
	return c.services.room
}

func (c *Container) ScheduleService(ctx context.Context) service.ScheduleService {
	if c.services.schedule == nil {
		c.services.schedule = schedule.NewScheduleService(c.ScheduleRepo(ctx), c.RoomRepo(ctx))
	}
	return c.services.schedule
}

func (c *Container) SlotService(ctx context.Context) service.SlotService {
	if c.services.slot == nil {
		c.services.slot = slot.NewSlotService(c.SlotRepo(ctx), c.ScheduleRepo(ctx), c.RoomRepo(ctx))
	}
	return c.services.slot
}
