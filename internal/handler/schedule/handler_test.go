package schedule

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"test-backend-1-kuprinvv/internal/mocks"
	"test-backend-1-kuprinvv/internal/model"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func validScheduleBody(roomID uuid.UUID) map[string]any {
	return map[string]any{
		"roomId":     roomID,
		"daysOfWeek": []int{1, 2, 3, 4, 5},
		"startTime":  "09:00",
		"endTime":    "18:00",
	}
}

func TestCreateSchedule(t *testing.T) {
	roomID := uuid.New()
	scheduleID := uuid.New()
	start := time.Date(1, 1, 1, 9, 0, 0, 0, time.UTC)
	end := time.Date(1, 1, 1, 18, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		roomID     string
		body       any
		setup      func(svc *mocks.MockScheduleService)
		wantStatus int
	}{
		{
			name:   "успешное создание",
			roomID: roomID.String(),
			body:   validScheduleBody(roomID),
			setup: func(svc *mocks.MockScheduleService) {
				svc.EXPECT().
					CreateSchedule(gomock.Any(), gomock.Any()).
					Return(&model.Schedule{ID: scheduleID, RoomID: roomID, StartTime: start, EndTime: end, DaysOfWeek: []int{1, 2, 3, 4, 5}}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "невалидный roomId",
			roomID:     "not-a-uuid",
			body:       validScheduleBody(roomID),
			setup:      func(svc *mocks.MockScheduleService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "невалидный JSON",
			roomID:     roomID.String(),
			body:       "not json",
			setup:      func(svc *mocks.MockScheduleService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "невалидный формат времени",
			roomID: roomID.String(),
			body: map[string]any{
				"daysOfWeek": []int{1},
				"startTime":  "not-a-time",
				"endTime":    "18:00",
			},
			setup:      func(svc *mocks.MockScheduleService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "endTime раньше startTime",
			roomID: roomID.String(),
			body: map[string]any{
				"daysOfWeek": []int{1},
				"startTime":  "18:00",
				"endTime":    "09:00",
			},
			setup:      func(svc *mocks.MockScheduleService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "переговорка не найдена",
			roomID: roomID.String(),
			body:   validScheduleBody(roomID),
			setup: func(svc *mocks.MockScheduleService) {
				svc.EXPECT().
					CreateSchedule(gomock.Any(), gomock.Any()).
					Return(nil, model.ErrRoomNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "расписание уже существует",
			roomID: roomID.String(),
			body:   validScheduleBody(roomID),
			setup: func(svc *mocks.MockScheduleService) {
				svc.EXPECT().
					CreateSchedule(gomock.Any(), gomock.Any()).
					Return(nil, model.ErrScheduleAlreadyExists)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:   "внутренняя ошибка",
			roomID: roomID.String(),
			body:   validScheduleBody(roomID),
			setup: func(svc *mocks.MockScheduleService) {
				svc.EXPECT().
					CreateSchedule(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockScheduleService(ctrl)
			tt.setup(svc)

			h := NewScheduleHandler(svc)

			rawBody, _ := json.Marshal(tt.body)
			r := httptest.NewRequest(http.MethodPost, "/rooms/"+tt.roomID+"/schedule/create", bytes.NewReader(rawBody))
			r.Header.Set("Content-Type", "application/json")
			r = withChiParam(r, roomIdParam, tt.roomID)
			w := httptest.NewRecorder()

			h.CreateSchedule(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
