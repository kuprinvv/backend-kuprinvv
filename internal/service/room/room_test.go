package room

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

func TestCreateRoom(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")

	cap10 := 10

	tests := []struct {
		name    string
		input   model.Room
		setup   func(repo *mocks.MockRoomRepository)
		wantErr error
	}{
		{
			name:  "успешное создание",
			input: model.Room{Name: "Room A", Capacity: &cap10},
			setup: func(repo *mocks.MockRoomRepository) {
				repo.EXPECT().
					CreateRoom(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, r model.Room) (*model.Room, error) {
						assert.NotEqual(t, uuid.Nil, r.ID, "ID must be assigned")
						assert.Equal(t, "Room A", r.Name)
						return &r, nil
					})
			},
		},
		{
			name:    "ошибка репозитория",
			input:   model.Room{Name: "Room B"},
			setup:   func(repo *mocks.MockRoomRepository) { repo.EXPECT().CreateRoom(ctx, gomock.Any()).Return(nil, dbErr) },
			wantErr: dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRoomRepository(ctrl)
			tt.setup(repo)

			svc := NewRoomService(repo)
			result, err := svc.CreateRoom(ctx, tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEqual(t, uuid.Nil, result.ID)
			}
		})
	}
}

func TestListRooms(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")

	rooms := []model.Room{
		{ID: uuid.New(), Name: "Room A"},
		{ID: uuid.New(), Name: "Room B"},
	}

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockRoomRepository)
		wantLen int
		wantErr error
	}{
		{
			name:    "возвращает переговорки",
			setup:   func(repo *mocks.MockRoomRepository) { repo.EXPECT().ListRooms(ctx).Return(rooms, nil) },
			wantLen: 2,
		},
		{
			name:    "возвращает пустой список",
			setup:   func(repo *mocks.MockRoomRepository) { repo.EXPECT().ListRooms(ctx).Return([]model.Room{}, nil) },
			wantLen: 0,
		},
		{
			name:    "ошибка репозитория",
			setup:   func(repo *mocks.MockRoomRepository) { repo.EXPECT().ListRooms(ctx).Return(nil, dbErr) },
			wantErr: dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRoomRepository(ctrl)
			tt.setup(repo)

			svc := NewRoomService(repo)
			result, err := svc.ListRooms(ctx)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}
