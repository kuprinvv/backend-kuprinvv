package di

import (
	"context"
	"test-backend-1-kuprinvv/internal/handler/auth"
	"test-backend-1-kuprinvv/internal/handler/booking"
	"test-backend-1-kuprinvv/internal/handler/room"
	"test-backend-1-kuprinvv/internal/handler/schedule"
	"test-backend-1-kuprinvv/internal/handler/slot"
)

type handlerProvider struct {
	auth     *auth.Handler
	booking  *booking.Handler
	room     *room.Handler
	schedule *schedule.Handler
	slot     *slot.Handler
}

func (c *Container) AuthHandler(ctx context.Context) *auth.Handler {
	if c.handlers.auth == nil {
		c.handlers.auth = auth.NewAuthHandler(c.AuthService(ctx))
	}
	return c.handlers.auth
}

func (c *Container) BookingHandler(ctx context.Context) *booking.Handler {
	if c.handlers.booking == nil {
		c.handlers.booking = booking.NewBookingHandler(c.BookingService(ctx))
	}
	return c.handlers.booking
}

func (c *Container) RoomHandler(ctx context.Context) *room.Handler {
	if c.handlers.room == nil {
		c.handlers.room = room.NewRoomHandler(c.RoomService(ctx))
	}
	return c.handlers.room
}

func (c *Container) ScheduleHandler(ctx context.Context) *schedule.Handler {
	if c.handlers.schedule == nil {
		c.handlers.schedule = schedule.NewScheduleHandler(c.ScheduleService(ctx))
	}
	return c.handlers.schedule
}

func (c *Container) SlotHandler(ctx context.Context) *slot.Handler {
	if c.handlers.slot == nil {
		c.handlers.slot = slot.NewSlotHandler(c.SlotService(ctx))
	}
	return c.handlers.slot
}
