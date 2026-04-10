package schedule

import "test-backend-1-kuprinvv/internal/service"

type Handler struct {
	scheduleServ service.ScheduleService
}

func NewScheduleHandler(scheduleServ service.ScheduleService) *Handler {
	return &Handler{scheduleServ}
}
