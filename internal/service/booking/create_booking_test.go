package booking

import (
	"context"
	"errors"
	"test-backend-1-kuprinvv/internal/mocks"
	"test-backend-1-kuprinvv/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateBooking(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	slotID := uuid.New()
	futureTime := time.Now().UTC().Add(24 * time.Hour)
	pastTime := time.Now().UTC().Add(-24 * time.Hour)
	confLink := "https://meet.example.com/abc"

	type fields struct {
		bookingRepo      *mocks.MockBookingRepository
		slotRepo         *mocks.MockSlotRepository
		conferenceClient *mocks.MockConferenceClient
	}

	tests := []struct {
		name                 string
		userID               uuid.UUID
		slotID               uuid.UUID
		createConferenceLink bool
		setup                func(f fields)
		wantErr              error
		checkResult          func(t *testing.T, b *model.Booking)
	}{
		{
			name:                 "успешное бронирование без конференц-ссылки",
			userID:               userID,
			slotID:               slotID,
			createConferenceLink: false,
			setup: func(f fields) {
				slot := &model.Slot{ID: slotID, RoomID: uuid.New(), StartTime: futureTime, EndTime: futureTime.Add(30 * time.Minute)}
				f.slotRepo.EXPECT().
					GetSlotByID(ctx, slotID).
					Return(slot, nil).
					Times(1)
				createdBooking := &model.Booking{ID: uuid.New(), UserID: userID, SlotID: slotID, Status: "active"}
				f.bookingRepo.EXPECT().
					CreateBooking(ctx, gomock.Any()).
					Return(createdBooking, nil).
					Times(1)
			},
			wantErr: nil,
			checkResult: func(t *testing.T, b *model.Booking) {
				require.NotNil(t, b)
				assert.Equal(t, userID, b.UserID)
				assert.Equal(t, slotID, b.SlotID)
				assert.Nil(t, b.ConferenceLink)
			},
		},
		{
			name:                 "успешное бронирование с конференц-ссылкой",
			userID:               userID,
			slotID:               slotID,
			createConferenceLink: true,
			setup: func(f fields) {
				slot := &model.Slot{ID: slotID, RoomID: uuid.New(), StartTime: futureTime, EndTime: futureTime.Add(30 * time.Minute)}
				f.slotRepo.EXPECT().
					GetSlotByID(ctx, slotID).
					Return(slot, nil).
					Times(1)
				f.conferenceClient.EXPECT().
					CreateLink(ctx).
					Return(confLink, nil).
					Times(1)
				createdBooking := &model.Booking{ID: uuid.New(), UserID: userID, SlotID: slotID, Status: "active", ConferenceLink: &confLink}
				f.bookingRepo.EXPECT().
					CreateBooking(ctx, gomock.Any()).
					Return(createdBooking, nil).
					Times(1)
			},
			wantErr: nil,
			checkResult: func(t *testing.T, b *model.Booking) {
				require.NotNil(t, b)
				assert.Equal(t, userID, b.UserID)
				assert.Equal(t, slotID, b.SlotID)
				require.NotNil(t, b.ConferenceLink)
				assert.Equal(t, confLink, *b.ConferenceLink)
			},
		},
		{
			name:                 "конференц-сервис недоступен — бронь создаётся без ссылки",
			userID:               userID,
			slotID:               slotID,
			createConferenceLink: true,
			setup: func(f fields) {
				slot := &model.Slot{ID: slotID, RoomID: uuid.New(), StartTime: futureTime, EndTime: futureTime.Add(30 * time.Minute)}
				f.slotRepo.EXPECT().
					GetSlotByID(ctx, slotID).
					Return(slot, nil).
					Times(1)
				f.conferenceClient.EXPECT().
					CreateLink(ctx).
					Return("", errors.New("conference service unavailable")).
					Times(1)
				createdBooking := &model.Booking{ID: uuid.New(), UserID: userID, SlotID: slotID, Status: "active"}
				f.bookingRepo.EXPECT().
					CreateBooking(ctx, gomock.Any()).
					Return(createdBooking, nil).
					Times(1)
			},
			wantErr: nil,
			checkResult: func(t *testing.T, b *model.Booking) {
				require.NotNil(t, b)
				assert.Nil(t, b.ConferenceLink)
			},
		},
		{
			name:                 "слот не найден",
			userID:               userID,
			slotID:               slotID,
			createConferenceLink: false,
			setup: func(f fields) {
				f.slotRepo.EXPECT().
					GetSlotByID(ctx, slotID).
					Return(nil, model.ErrSlotNotFound).
					Times(1)
			},
			wantErr: model.ErrSlotNotFound,
			checkResult: func(t *testing.T, b *model.Booking) {
				assert.Nil(t, b)
			},
		},
		{
			name:                 "слот в прошлом",
			userID:               userID,
			slotID:               slotID,
			createConferenceLink: false,
			setup: func(f fields) {
				slot := &model.Slot{ID: slotID, RoomID: uuid.New(), StartTime: pastTime, EndTime: pastTime.Add(30 * time.Minute)}
				f.slotRepo.EXPECT().
					GetSlotByID(ctx, slotID).
					Return(slot, nil).
					Times(1)
			},
			wantErr: model.ErrPastSlot,
			checkResult: func(t *testing.T, b *model.Booking) {
				assert.Nil(t, b)
			},
		},
		{
			name:                 "слот уже занят — нарушение уникальности",
			userID:               userID,
			slotID:               slotID,
			createConferenceLink: false,
			setup: func(f fields) {
				slot := &model.Slot{ID: slotID, RoomID: uuid.New(), StartTime: futureTime, EndTime: futureTime.Add(30 * time.Minute)}
				f.slotRepo.EXPECT().
					GetSlotByID(ctx, slotID).
					Return(slot, nil).
					Times(1)
				pgErr := &pgconn.PgError{Code: "23505"}
				f.bookingRepo.EXPECT().
					CreateBooking(ctx, gomock.Any()).
					Return(nil, pgErr).
					Times(1)
			},
			wantErr: model.ErrSlotAlreadyBooked,
			checkResult: func(t *testing.T, b *model.Booking) {
				assert.Nil(t, b)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				bookingRepo:      mocks.NewMockBookingRepository(ctrl),
				slotRepo:         mocks.NewMockSlotRepository(ctrl),
				conferenceClient: mocks.NewMockConferenceClient(ctrl),
			}
			tt.setup(f)

			svc := NewBookingService(f.bookingRepo, f.slotRepo, f.conferenceClient)
			got, err := svc.CreateBooking(ctx, tt.userID, tt.slotID, tt.createConferenceLink)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			tt.checkResult(t, got)
		})
	}
}
