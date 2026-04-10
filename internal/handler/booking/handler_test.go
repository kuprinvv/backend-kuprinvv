package booking

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"test-backend-1-kuprinvv/internal/middleware"
	"test-backend-1-kuprinvv/internal/mocks"
	"test-backend-1-kuprinvv/internal/model"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// withChiParam добавляет URL-параметр chi в контекст запроса.
func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// withUser добавляет пользователя в контекст запроса.
func withUser(r *http.Request, userID uuid.UUID, role string) *http.Request {
	ctx := middleware.WithUser(r.Context(), middleware.UserContext{UserID: userID, Role: role})
	return r.WithContext(ctx)
}

func TestCreateBooking(t *testing.T) {
	userID := uuid.New()
	slotID := uuid.New()
	bookingID := uuid.New()
	now := time.Now()

	tests := []struct {
		name       string
		body       any
		setup      func(svc *mocks.MockBookingService)
		wantStatus int
	}{
		{
			name: "успешное создание",
			body: map[string]any{"slotId": slotID},
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().
					CreateBooking(gomock.Any(), userID, slotID, false).
					Return(&model.Booking{ID: bookingID, UserID: userID, SlotID: slotID, Status: "active", CreatedAt: &now}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "невалидный JSON",
			body:       "not json",
			setup:      func(svc *mocks.MockBookingService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "слот не найден",
			body: map[string]any{"slotId": slotID},
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().
					CreateBooking(gomock.Any(), userID, slotID, false).
					Return(nil, model.ErrSlotNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "слот уже занят",
			body: map[string]any{"slotId": slotID},
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().
					CreateBooking(gomock.Any(), userID, slotID, false).
					Return(nil, model.ErrSlotAlreadyBooked)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "слот в прошлом",
			body: map[string]any{"slotId": slotID},
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().
					CreateBooking(gomock.Any(), userID, slotID, false).
					Return(nil, model.ErrPastSlot)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "внутренняя ошибка",
			body: map[string]any{"slotId": slotID},
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().
					CreateBooking(gomock.Any(), userID, slotID, false).
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockBookingService(ctrl)
			tt.setup(svc)

			h := NewBookingHandler(svc)

			rawBody, _ := json.Marshal(tt.body)
			r := httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewReader(rawBody))
			r.Header.Set("Content-Type", "application/json")
			r = withUser(r, userID, "user")
			w := httptest.NewRecorder()

			h.CreateBooking(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCancelBooking(t *testing.T) {
	userID := uuid.New()
	bookingID := uuid.New()
	now := time.Now()

	tests := []struct {
		name       string
		bookingID  string
		setup      func(svc *mocks.MockBookingService)
		wantStatus int
	}{
		{
			name:      "успешная отмена",
			bookingID: bookingID.String(),
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().
					CancelBooking(gomock.Any(), userID, bookingID).
					Return(&model.Booking{ID: bookingID, UserID: userID, Status: "cancelled", CreatedAt: &now}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "невалидный UUID",
			bookingID:  "not-a-uuid",
			setup:      func(svc *mocks.MockBookingService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:      "чужая бронь",
			bookingID: bookingID.String(),
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().
					CancelBooking(gomock.Any(), userID, bookingID).
					Return(nil, model.ErrForbidden)
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:      "бронь не найдена",
			bookingID: bookingID.String(),
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().
					CancelBooking(gomock.Any(), userID, bookingID).
					Return(nil, model.ErrBookingNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:      "внутренняя ошибка",
			bookingID: bookingID.String(),
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().
					CancelBooking(gomock.Any(), userID, bookingID).
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockBookingService(ctrl)
			tt.setup(svc)

			h := NewBookingHandler(svc)

			r := httptest.NewRequest(http.MethodPost, "/bookings/"+tt.bookingID+"/cancel", nil)
			r = withUser(r, userID, "user")
			r = withChiParam(r, bookingIdURLParam, tt.bookingID)
			w := httptest.NewRecorder()

			h.CancelBooking(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestListBookings(t *testing.T) {
	now := time.Now()
	bookings := []model.Booking{
		{ID: uuid.New(), Status: "active", CreatedAt: &now},
	}
	pagination := model.Pagination{Page: 1, PageSize: 20, Total: 1}

	tests := []struct {
		name       string
		query      string
		setup      func(svc *mocks.MockBookingService)
		wantStatus int
	}{
		{
			name:  "успешно, дефолтные параметры",
			query: "",
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().ListBookings(gomock.Any(), 1, 20).Return(bookings, &pagination, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "кастомная пагинация",
			query: "?page=2&pageSize=10",
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().ListBookings(gomock.Any(), 2, 10).Return(bookings, &pagination, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "невалидный page",
			query:      "?page=abc",
			setup:      func(svc *mocks.MockBookingService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "страница меньше 1",
			query:      "?page=0",
			setup:      func(svc *mocks.MockBookingService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "размер страницы больше 100",
			query:      "?pageSize=101",
			setup:      func(svc *mocks.MockBookingService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "внутренняя ошибка",
			query: "",
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().ListBookings(gomock.Any(), 1, 20).Return(nil, nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockBookingService(ctrl)
			tt.setup(svc)

			h := NewBookingHandler(svc)

			r := httptest.NewRequest(http.MethodGet, "/bookings/list"+tt.query, nil)
			w := httptest.NewRecorder()

			h.ListBookings(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetMyBookings(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name       string
		setup      func(svc *mocks.MockBookingService)
		wantStatus int
	}{
		{
			name: "успешно",
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().
					GetMyBookings(gomock.Any(), userID).
					Return([]model.Booking{{ID: uuid.New(), Status: "active", CreatedAt: &now}}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "пустой список",
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().GetMyBookings(gomock.Any(), userID).Return([]model.Booking{}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "внутренняя ошибка",
			setup: func(svc *mocks.MockBookingService) {
				svc.EXPECT().GetMyBookings(gomock.Any(), userID).Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockBookingService(ctrl)
			tt.setup(svc)

			h := NewBookingHandler(svc)

			r := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)
			r = withUser(r, userID, "user")
			w := httptest.NewRecorder()

			h.GetMyBookings(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestListBookings_ResponseBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockBookingService(ctrl)

	now := time.Now()
	bookingID := uuid.New()
	svc.EXPECT().ListBookings(gomock.Any(), 1, 20).Return(
		[]model.Booking{{ID: bookingID, Status: "active", CreatedAt: &now}},
		&model.Pagination{Page: 1, PageSize: 20, Total: 1},
		nil,
	)

	h := NewBookingHandler(svc)
	r := httptest.NewRequest(http.MethodGet, "/bookings/list", nil)
	w := httptest.NewRecorder()
	h.ListBookings(w, r)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Bookings   []any          `json:"bookings"`
		Pagination map[string]any `json:"pagination"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Bookings, 1)
	assert.Equal(t, float64(1), resp.Pagination["total"])
}
