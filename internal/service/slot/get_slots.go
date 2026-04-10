package slot

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"
	"time"

	"github.com/google/uuid"
)

func (s *serv) GetSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]model.Slot, error) {
	slots, err := s.slotRepo.GetSlotsByRoomAndDate(ctx, roomID, date)
	if err != nil {
		return nil, err
	}

	if len(slots) > 0 {
		return slots, nil
	}

	_, err = s.roomRepo.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	schedule, err := s.scheduleRepo.GetByRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}

	if schedule == nil {
		return []model.Slot{}, nil
	}

	if !isDateAllowed(date, schedule.DaysOfWeek) {
		return []model.Slot{}, nil
	}

	err = s.generateSlots(ctx, roomID, date, schedule)
	if err != nil {
		return nil, err
	}

	return s.slotRepo.GetSlotsByRoomAndDate(ctx, roomID, date)
}

func isDateAllowed(date time.Time, days []int) bool {
	weekday := int(date.Weekday()) // 0–6

	apiDay := weekday
	if apiDay == 0 {
		apiDay = 7
	}

	for _, d := range days {
		if d == apiDay {
			return true
		}
	}

	return false
}

func (s *serv) generateSlots(
	ctx context.Context,
	roomID uuid.UUID,
	date time.Time,
	schedule *model.Schedule,
) error {
	dayStart := time.Date(
		date.Year(), date.Month(), date.Day(),
		0, 0, 0, 0, time.UTC,
	)

	start := time.Date(
		dayStart.Year(), dayStart.Month(), dayStart.Day(),
		schedule.StartTime.Hour(),
		schedule.StartTime.Minute(),
		0, 0,
		time.UTC,
	)

	end := time.Date(
		dayStart.Year(), dayStart.Month(), dayStart.Day(),
		schedule.EndTime.Hour(),
		schedule.EndTime.Minute(),
		0, 0,
		time.UTC,
	)

	var slots []model.Slot
	current := start

	for current.Add(30*time.Minute).Before(end) || current.Add(30*time.Minute).Equal(end) {

		slotEnd := current.Add(30 * time.Minute)

		if slotEnd.Before(time.Now().UTC()) {
			current = slotEnd
			continue
		}

		slots = append(slots, model.Slot{
			ID:        uuid.New(),
			RoomID:    roomID,
			StartTime: current,
			EndTime:   slotEnd,
		})

		current = slotEnd
	}

	if len(slots) == 0 {
		return nil
	}

	return s.slotRepo.CreateSlots(ctx, slots)
}
