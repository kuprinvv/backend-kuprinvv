package di

import (
	"context"
	"test-backend-1-kuprinvv/internal/repository"
	"test-backend-1-kuprinvv/internal/repository/booking"
	"test-backend-1-kuprinvv/internal/repository/room"
	"test-backend-1-kuprinvv/internal/repository/schedule"
	"test-backend-1-kuprinvv/internal/repository/slot"
	"test-backend-1-kuprinvv/internal/repository/user"
)

type repositoryProvider struct {
	booking  repository.BookingRepository
	room     repository.RoomRepository
	schedule repository.ScheduleRepository
	slot     repository.SlotRepository
	user     repository.UserRepository
}

func (c *Container) BookingRepo(ctx context.Context) repository.BookingRepository {
	if c.repos.booking == nil {
		c.repos.booking = booking.NewBookingRepository(c.DB(ctx))
	}
	return c.repos.booking
}

func (c *Container) RoomRepo(ctx context.Context) repository.RoomRepository {
	if c.repos.room == nil {
		c.repos.room = room.NewRoomRepository(c.DB(ctx))
	}
	return c.repos.room
}

func (c *Container) ScheduleRepo(ctx context.Context) repository.ScheduleRepository {
	if c.repos.schedule == nil {
		c.repos.schedule = schedule.NewUserRepository(c.DB(ctx))
	}
	return c.repos.schedule
}

func (c *Container) SlotRepo(ctx context.Context) repository.SlotRepository {
	if c.repos.slot == nil {
		c.repos.slot = slot.NewSlotRepository(c.DB(ctx))
	}
	return c.repos.slot
}

func (c *Container) UserRepo(ctx context.Context) repository.UserRepository {
	if c.repos.user == nil {
		c.repos.user = user.NewUserRepository(c.DB(ctx))
	}
	return c.repos.user
}
