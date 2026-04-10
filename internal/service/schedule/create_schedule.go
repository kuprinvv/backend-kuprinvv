package schedule

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/internal/service"

	"github.com/google/uuid"
)

func (s *serv) CreateSchedule(ctx context.Context, schedule model.Schedule) (*model.Schedule, error) {
	schedule.ID = uuid.New()

	_, err := s.roomRepo.GetRoomByID(ctx, schedule.RoomID)
	if err != nil {
		return nil, err
	}

	err = s.scheduleRepo.CreateSchedule(ctx, schedule)
	if err != nil {
		if service.IsUniqueViolation(err) {
			return nil, model.ErrScheduleAlreadyExists
		}
		return nil, err
	}

	return &schedule, nil
}
