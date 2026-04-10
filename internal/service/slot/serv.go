package slot

import (
	"test-backend-1-kuprinvv/internal/repository"
	"test-backend-1-kuprinvv/internal/service"
)

var _ service.SlotService = (*serv)(nil)

type serv struct {
	slotRepo     repository.SlotRepository
	scheduleRepo repository.ScheduleRepository
	roomRepo     repository.RoomRepository
}

func NewSlotService(
	slotRepo repository.SlotRepository,
	scheduleRepo repository.ScheduleRepository,
	roomRepo repository.RoomRepository,
) *serv {
	return &serv{
		slotRepo:     slotRepo,
		scheduleRepo: scheduleRepo,
		roomRepo:     roomRepo,
	}
}
