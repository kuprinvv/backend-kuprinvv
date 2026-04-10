package booking

import (
	"context"
	"test-backend-1-kuprinvv/internal/mocks"
	"test-backend-1-kuprinvv/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCancelBooking(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()
	otherUserID := uuid.New()
	bookingID := uuid.New()

	type fields struct {
		bookingRepo      *mocks.MockBookingRepository
		slotRepo         *mocks.MockSlotRepository
		conferenceClient *mocks.MockConferenceClient
	}

	tests := []struct {
		name        string
		userID      uuid.UUID
		bookingID   uuid.UUID
		setup       func(f fields)
		wantErr     error
		checkResult func(t *testing.T, b *model.Booking)
	}{
		{
			name:      "успешная отмена активной брони владельцем",
			userID:    ownerID,
			bookingID: bookingID,
			setup: func(f fields) {
				existing := &model.Booking{ID: bookingID, UserID: ownerID, Status: "active"}
				f.bookingRepo.EXPECT().
					GetBookingByID(ctx, bookingID).
					Return(existing, nil).
					Times(1)
				cancelled := &model.Booking{ID: bookingID, UserID: ownerID, Status: "cancelled"}
				f.bookingRepo.EXPECT().
					CancelBooking(ctx, bookingID).
					Return(cancelled, nil).
					Times(1)
			},
			wantErr: nil,
			checkResult: func(t *testing.T, b *model.Booking) {
				require.NotNil(t, b)
				assert.Equal(t, "cancelled", b.Status)
				assert.Equal(t, bookingID, b.ID)
			},
		},
		{
			name:      "идемпотентность: уже отменена — повторного вызова репозитория нет",
			userID:    ownerID,
			bookingID: bookingID,
			setup: func(f fields) {
				existing := &model.Booking{ID: bookingID, UserID: ownerID, Status: "cancelled"}
				f.bookingRepo.EXPECT().
					GetBookingByID(ctx, bookingID).
					Return(existing, nil).
					Times(1)
				// CancelBooking must NOT be called again
			},
			wantErr: nil,
			checkResult: func(t *testing.T, b *model.Booking) {
				require.NotNil(t, b)
				assert.Equal(t, "cancelled", b.Status)
			},
		},
		{
			name:      "бронь не найдена",
			userID:    ownerID,
			bookingID: bookingID,
			setup: func(f fields) {
				f.bookingRepo.EXPECT().
					GetBookingByID(ctx, bookingID).
					Return(nil, model.ErrBookingNotFound).
					Times(1)
			},
			wantErr: model.ErrBookingNotFound,
			checkResult: func(t *testing.T, b *model.Booking) {
				assert.Nil(t, b)
			},
		},
		{
			name:      "запрет: другой пользователь",
			userID:    otherUserID,
			bookingID: bookingID,
			setup: func(f fields) {
				existing := &model.Booking{ID: bookingID, UserID: ownerID, Status: "active"}
				f.bookingRepo.EXPECT().
					GetBookingByID(ctx, bookingID).
					Return(existing, nil).
					Times(1)
			},
			wantErr: model.ErrForbidden,
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
			got, err := svc.CancelBooking(ctx, tt.userID, tt.bookingID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			tt.checkResult(t, got)
		})
	}
}
