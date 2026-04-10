package schedule

import (
	"context"
	"errors"
	"test-backend-1-kuprinvv/internal/mocks"
	"test-backend-1-kuprinvv/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateSchedule(t *testing.T) {
	ctx := context.Background()
	roomID := uuid.New()
	room := &model.Room{ID: roomID, Name: "Room A"}
	dbErr := errors.New("db error")

	tests := []struct {
		name    string
		input   model.Schedule
		setup   func(scheduleRepo *mocks.MockScheduleRepository, roomRepo *mocks.MockRoomRepository)
		wantErr bool
	}{
		{
			name:  "успешное создание",
			input: model.Schedule{RoomID: roomID, DaysOfWeek: []int{1, 2, 3}},
			setup: func(scheduleRepo *mocks.MockScheduleRepository, roomRepo *mocks.MockRoomRepository) {
				roomRepo.EXPECT().GetRoomByID(ctx, roomID).Return(room, nil)
				scheduleRepo.EXPECT().
					CreateSchedule(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, s model.Schedule) error {
						assert.NotEqual(t, uuid.Nil, s.ID, "ID must be assigned")
						return nil
					})
			},
		},
		{
			name:  "переговорка не найдена",
			input: model.Schedule{RoomID: roomID},
			setup: func(scheduleRepo *mocks.MockScheduleRepository, roomRepo *mocks.MockRoomRepository) {
				roomRepo.EXPECT().GetRoomByID(ctx, roomID).Return(nil, model.ErrRoomNotFound)
			},
			wantErr: true,
		},
		{
			name:  "ошибка репозитория расписаний",
			input: model.Schedule{RoomID: roomID},
			setup: func(scheduleRepo *mocks.MockScheduleRepository, roomRepo *mocks.MockRoomRepository) {
				roomRepo.EXPECT().GetRoomByID(ctx, roomID).Return(room, nil)
				scheduleRepo.EXPECT().CreateSchedule(ctx, gomock.Any()).Return(dbErr)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			scheduleRepo := mocks.NewMockScheduleRepository(ctrl)
			roomRepo := mocks.NewMockRoomRepository(ctrl)
			tt.setup(scheduleRepo, roomRepo)

			svc := NewScheduleService(scheduleRepo, roomRepo)
			result, err := svc.CreateSchedule(ctx, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEqual(t, uuid.Nil, result.ID)
				assert.Equal(t, roomID, result.RoomID)
			}
		})
	}
}
