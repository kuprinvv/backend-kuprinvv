package slot

import (
	"context"
	"test-backend-1-kuprinvv/internal/mocks"
	"test-backend-1-kuprinvv/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetSlots(t *testing.T) {
	ctx := context.Background()
	roomID := uuid.New()

	// Monday far in the future so all generated slots pass the "not in the past" guard.
	// time.Weekday: Monday == 1; API day: 1 = Mon.
	monday := nextWeekday(time.Now().UTC().AddDate(0, 0, 7), time.Monday)
	date := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)

	sunday := nextWeekday(date.AddDate(0, 0, 1), time.Sunday)
	sundayDate := time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 0, 0, 0, 0, time.UTC)

	existingSlots := []model.Slot{
		{ID: uuid.New(), RoomID: roomID, StartTime: date.Add(9 * time.Hour), EndTime: date.Add(9*time.Hour + 30*time.Minute)},
	}

	schedule := &model.Schedule{
		ID:         uuid.New(),
		RoomID:     roomID,
		DaysOfWeek: []int{1, 2, 3, 4, 5}, // Mon–Fri
		StartTime:  time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC),
		EndTime:    time.Date(2000, 1, 1, 10, 0, 0, 0, time.UTC),
	}

	type fields struct {
		slotRepo     *mocks.MockSlotRepository
		scheduleRepo *mocks.MockScheduleRepository
		roomRepo     *mocks.MockRoomRepository
	}

	tests := []struct {
		name        string
		roomID      uuid.UUID
		date        time.Time
		setup       func(f fields)
		wantErr     bool
		checkResult func(t *testing.T, slots []model.Slot)
	}{
		{
			name:   "слоты уже есть в БД — возвращаем сразу",
			roomID: roomID,
			date:   date,
			setup: func(f fields) {
				f.slotRepo.EXPECT().
					GetSlotsByRoomAndDate(ctx, roomID, date).
					Return(existingSlots, nil).
					Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, slots []model.Slot) {
				require.Len(t, slots, 1)
				assert.Equal(t, existingSlots[0].ID, slots[0].ID)
			},
		},
		{
			name:   "слотов нет — генерация: комната и расписание есть, дата допустима (понедельник)",
			roomID: roomID,
			date:   date,
			setup: func(f fields) {
				generatedSlots := []model.Slot{
					{ID: uuid.New(), RoomID: roomID, StartTime: date.Add(9 * time.Hour), EndTime: date.Add(9*time.Hour + 30*time.Minute)},
					{ID: uuid.New(), RoomID: roomID, StartTime: date.Add(9*time.Hour + 30*time.Minute), EndTime: date.Add(10 * time.Hour)},
				}

				firstCall := f.slotRepo.EXPECT().
					GetSlotsByRoomAndDate(ctx, roomID, date).
					Return([]model.Slot{}, nil)
				secondCall := f.slotRepo.EXPECT().
					GetSlotsByRoomAndDate(ctx, roomID, date).
					Return(generatedSlots, nil)
				gomock.InOrder(firstCall, secondCall)

				f.roomRepo.EXPECT().
					GetRoomByID(ctx, roomID).
					Return(&model.Room{ID: roomID}, nil).
					Times(1)
				f.scheduleRepo.EXPECT().
					GetByRoom(ctx, roomID).
					Return(schedule, nil).
					Times(1)
				f.slotRepo.EXPECT().
					CreateSlots(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, slots []model.Slot) {
				assert.NotEmpty(t, slots)
			},
		},
		{
			name:   "слотов нет, расписания нет — возвращаем пустой список",
			roomID: roomID,
			date:   date,
			setup: func(f fields) {
				f.slotRepo.EXPECT().
					GetSlotsByRoomAndDate(ctx, roomID, date).
					Return([]model.Slot{}, nil).
					Times(1)
				f.roomRepo.EXPECT().
					GetRoomByID(ctx, roomID).
					Return(&model.Room{ID: roomID}, nil).
					Times(1)
				f.scheduleRepo.EXPECT().
					GetByRoom(ctx, roomID).
					Return(nil, nil).
					Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, slots []model.Slot) {
				assert.Empty(t, slots)
			},
		},
		{
			name:   "слотов нет, дата недопустима (воскресенье, расписание только в будни)",
			roomID: roomID,
			date:   sundayDate,
			setup: func(f fields) {
				f.slotRepo.EXPECT().
					GetSlotsByRoomAndDate(ctx, roomID, sundayDate).
					Return([]model.Slot{}, nil).
					Times(1)
				f.roomRepo.EXPECT().
					GetRoomByID(ctx, roomID).
					Return(&model.Room{ID: roomID}, nil).
					Times(1)
				f.scheduleRepo.EXPECT().
					GetByRoom(ctx, roomID).
					Return(schedule, nil).
					Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, slots []model.Slot) {
				assert.Empty(t, slots)
			},
		},
		{
			name:   "переговорка не найдена — ошибка",
			roomID: roomID,
			date:   date,
			setup: func(f fields) {
				f.slotRepo.EXPECT().
					GetSlotsByRoomAndDate(ctx, roomID, date).
					Return([]model.Slot{}, nil).
					Times(1)
				f.roomRepo.EXPECT().
					GetRoomByID(ctx, roomID).
					Return(nil, model.ErrRoomNotFound).
					Times(1)
			},
			wantErr: true,
			checkResult: func(t *testing.T, slots []model.Slot) {
				assert.Nil(t, slots)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				slotRepo:     mocks.NewMockSlotRepository(ctrl),
				scheduleRepo: mocks.NewMockScheduleRepository(ctrl),
				roomRepo:     mocks.NewMockRoomRepository(ctrl),
			}
			tt.setup(f)

			svc := NewSlotService(f.slotRepo, f.scheduleRepo, f.roomRepo)
			got, err := svc.GetSlots(ctx, tt.roomID, tt.date)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			tt.checkResult(t, got)
		})
	}
}

// nextWeekday returns the next occurrence of the given weekday on or after start.
func nextWeekday(start time.Time, weekday time.Weekday) time.Time {
	d := start
	for d.Weekday() != weekday {
		d = d.AddDate(0, 0, 1)
	}
	return d
}
