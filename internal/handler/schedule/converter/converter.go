package converter

import (
	"fmt"
	"test-backend-1-kuprinvv/internal/handler/schedule/dto"
	"test-backend-1-kuprinvv/internal/model"
	"time"
)

func CreateScheduleCreateToService(req dto.CreateScheduleRequest) (model.Schedule, error) {
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return model.Schedule{}, fmt.Errorf("invalid startTime: must be HH:MM")
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return model.Schedule{}, fmt.Errorf("invalid endTime: must be HH:MM")
	}

	if !endTime.After(startTime) {
		return model.Schedule{}, fmt.Errorf("endTime must be after startTime")
	}

	return model.Schedule{
		ID:         req.ID,
		RoomID:     req.RoomID,
		DaysOfWeek: req.DaysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
	}, nil
}

func ServiceToCreateScheduleResponse(schedule model.Schedule) dto.CreateScheduleResponse {
	return dto.CreateScheduleResponse{
		Schedule: dto.Schedule{
			ID:         schedule.ID,
			RoomID:     schedule.RoomID,
			DaysOfWeek: schedule.DaysOfWeek,
			StartTime:  schedule.StartTime.UTC().Format("15:04"),
			EndTime:    schedule.EndTime.UTC().Format("15:04"),
		},
	}
}
