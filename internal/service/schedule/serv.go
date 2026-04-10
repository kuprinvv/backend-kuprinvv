package schedule

import (
	"test-backend-1-kuprinvv/internal/repository"
	"test-backend-1-kuprinvv/internal/service"
)

var _ service.ScheduleService = (*serv)(nil)

type serv struct {
	scheduleRepo repository.ScheduleRepository
	roomRepo     repository.RoomRepository
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, room repository.RoomRepository) *serv {
	return &serv{
		scheduleRepo: scheduleRepo,
		roomRepo:     room,
	}
}
