package dto

import (
	"github.com/google/uuid"
)

type CreateScheduleRequest struct {
	ID         uuid.UUID `json:"id"`
	RoomID     uuid.UUID `json:"roomId"`
	DaysOfWeek []int     `json:"daysOfWeek" validate:"min=1,dive,min=1,max=7"`
	StartTime  string    `json:"startTime" validate:"required"`
	EndTime    string    `json:"endTime" validate:"required"`
}

type CreateScheduleResponse struct {
	Schedule Schedule `json:"schedule"`
}

type Schedule struct {
	ID         uuid.UUID `json:"id"`
	RoomID     uuid.UUID `json:"roomId"`
	DaysOfWeek []int     `json:"daysOfWeek"`
	StartTime  string    `json:"startTime"`
	EndTime    string    `json:"endTime"`
}
