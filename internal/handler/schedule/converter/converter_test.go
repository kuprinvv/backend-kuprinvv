package converter_test

import (
	"test-backend-1-kuprinvv/internal/handler/schedule/converter"
	"test-backend-1-kuprinvv/internal/handler/schedule/dto"
	"test-backend-1-kuprinvv/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateScheduleCreateToService(t *testing.T) {
	roomID := uuid.New()

	req := dto.CreateScheduleRequest{
		RoomID:     roomID,
		DaysOfWeek: []int{1, 2, 3},
		StartTime:  "09:00",
		EndTime:    "18:00",
	}

	s, err := converter.CreateScheduleCreateToService(req)
	require.NoError(t, err)
	assert.Equal(t, roomID, s.RoomID)
	assert.Equal(t, []int{1, 2, 3}, s.DaysOfWeek)
	assert.Equal(t, 9, s.StartTime.Hour())
	assert.Equal(t, 0, s.StartTime.Minute())
	assert.Equal(t, 18, s.EndTime.Hour())
	assert.Equal(t, 0, s.EndTime.Minute())
}

func TestCreateScheduleCreateToService_InvalidFormat(t *testing.T) {
	req := dto.CreateScheduleRequest{
		DaysOfWeek: []int{1},
		StartTime:  "not-a-time",
		EndTime:    "18:00",
	}
	_, err := converter.CreateScheduleCreateToService(req)
	assert.Error(t, err)
}

func TestCreateScheduleCreateToService_EndBeforeStart(t *testing.T) {
	req := dto.CreateScheduleRequest{
		DaysOfWeek: []int{1},
		StartTime:  "18:00",
		EndTime:    "09:00",
	}
	_, err := converter.CreateScheduleCreateToService(req)
	assert.Error(t, err)
}

func TestServiceToCreateScheduleResponse(t *testing.T) {
	roomID := uuid.New()
	schedID := uuid.New()
	start := time.Date(1, 1, 1, 9, 0, 0, 0, time.UTC)
	end := time.Date(1, 1, 1, 18, 0, 0, 0, time.UTC)

	sched := model.Schedule{
		ID:         schedID,
		RoomID:     roomID,
		DaysOfWeek: []int{1, 5},
		StartTime:  start,
		EndTime:    end,
	}

	resp := converter.ServiceToCreateScheduleResponse(sched)
	assert.Equal(t, schedID, resp.Schedule.ID)
	assert.Equal(t, roomID, resp.Schedule.RoomID)
	assert.Equal(t, []int{1, 5}, resp.Schedule.DaysOfWeek)
	assert.Equal(t, "09:00", resp.Schedule.StartTime)
	assert.Equal(t, "18:00", resp.Schedule.EndTime)
}
