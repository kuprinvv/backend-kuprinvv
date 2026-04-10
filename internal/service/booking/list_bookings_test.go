package booking

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

func TestListBookings(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")

	bookings := []model.Booking{
		{ID: uuid.New(), Status: "active"},
		{ID: uuid.New(), Status: "cancelled"},
	}

	tests := []struct {
		name        string
		page        int
		pageSize    int
		setup       func(repo *mocks.MockBookingRepository)
		wantErr     error
		checkResult func(t *testing.T, bookings []model.Booking, p *model.Pagination)
	}{
		{
			name:     "успешно с пагинацией",
			page:     2,
			pageSize: 10,
			setup: func(repo *mocks.MockBookingRepository) {
				// offset = (2-1)*10 = 10
				repo.EXPECT().ListBookings(ctx, 10, 10).Return(bookings, 25, nil)
			},
			checkResult: func(t *testing.T, b []model.Booking, p *model.Pagination) {
				assert.Len(t, b, 2)
				require.NotNil(t, p)
				assert.Equal(t, 25, p.Total)
				assert.Equal(t, 2, p.Page)
				assert.Equal(t, 10, p.PageSize)
			},
		},
		{
			name:     "первая страница",
			page:     1,
			pageSize: 5,
			setup: func(repo *mocks.MockBookingRepository) {
				// offset = (1-1)*5 = 0
				repo.EXPECT().ListBookings(ctx, 5, 0).Return([]model.Booking{}, 0, nil)
			},
			checkResult: func(t *testing.T, b []model.Booking, p *model.Pagination) {
				assert.Empty(t, b)
				assert.Equal(t, 0, p.Total)
			},
		},
		{
			name:     "ошибка репозитория",
			page:     1,
			pageSize: 10,
			setup: func(repo *mocks.MockBookingRepository) {
				repo.EXPECT().ListBookings(ctx, 10, 0).Return(nil, 0, dbErr)
			},
			wantErr: dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bookingRepo := mocks.NewMockBookingRepository(ctrl)
			slotRepo := mocks.NewMockSlotRepository(ctrl)
			conferenceClient := mocks.NewMockConferenceClient(ctrl)
			tt.setup(bookingRepo)

			svc := NewBookingService(bookingRepo, slotRepo, conferenceClient)
			result, pagination, err := svc.ListBookings(ctx, tt.page, tt.pageSize)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
				assert.Nil(t, pagination)
			} else {
				require.NoError(t, err)
				tt.checkResult(t, result, pagination)
			}
		})
	}
}

func TestGetMyBookings(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	dbErr := errors.New("db error")

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockBookingRepository)
		wantLen int
		wantErr error
	}{
		{
			name: "возвращает будущие брони",
			setup: func(repo *mocks.MockBookingRepository) {
				repo.EXPECT().GetFutureBookingsByUser(ctx, userID).Return([]model.Booking{
					{ID: uuid.New(), UserID: userID, Status: "active"},
				}, nil)
			},
			wantLen: 1,
		},
		{
			name: "возвращает пустой список",
			setup: func(repo *mocks.MockBookingRepository) {
				repo.EXPECT().GetFutureBookingsByUser(ctx, userID).Return([]model.Booking{}, nil)
			},
			wantLen: 0,
		},
		{
			name: "ошибка репозитория",
			setup: func(repo *mocks.MockBookingRepository) {
				repo.EXPECT().GetFutureBookingsByUser(ctx, userID).Return(nil, dbErr)
			},
			wantErr: dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bookingRepo := mocks.NewMockBookingRepository(ctrl)
			slotRepo := mocks.NewMockSlotRepository(ctrl)
			conferenceClient := mocks.NewMockConferenceClient(ctrl)
			tt.setup(bookingRepo)

			svc := NewBookingService(bookingRepo, slotRepo, conferenceClient)
			result, err := svc.GetMyBookings(ctx, userID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}
