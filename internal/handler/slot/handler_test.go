package slot

import (
	"context"
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

func TestGetSlots(t *testing.T) {
	roomID := uuid.New()
	tomorrow := time.Now().UTC().Add(24 * time.Hour)

	tests := []struct {
		name       string
		roomID     string
		date       string
		setup      func(svc *mocks.MockSlotService)
		wantStatus int
	}{
		{
			name:   "успешно",
			roomID: roomID.String(),
			date:   tomorrow.Format(time.DateOnly),
			setup: func(svc *mocks.MockSlotService) {
				svc.EXPECT().
					GetSlots(gomock.Any(), roomID, gomock.Any()).
					Return([]model.Slot{{ID: uuid.New(), RoomID: roomID, StartTime: tomorrow, EndTime: tomorrow.Add(30 * time.Minute)}}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "невалидный roomId",
			roomID:     "not-a-uuid",
			date:       tomorrow.Format(time.DateOnly),
			setup:      func(svc *mocks.MockSlotService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "отсутствует date",
			roomID:     roomID.String(),
			date:       "",
			setup:      func(svc *mocks.MockSlotService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "невалидный формат date",
			roomID:     roomID.String(),
			date:       "25-03-2026",
			setup:      func(svc *mocks.MockSlotService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "переговорка не найдена",
			roomID: roomID.String(),
			date:   tomorrow.Format(time.DateOnly),
			setup: func(svc *mocks.MockSlotService) {
				svc.EXPECT().
					GetSlots(gomock.Any(), roomID, gomock.Any()).
					Return(nil, model.ErrRoomNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "внутренняя ошибка",
			roomID: roomID.String(),
			date:   tomorrow.Format(time.DateOnly),
			setup: func(svc *mocks.MockSlotService) {
				svc.EXPECT().
					GetSlots(gomock.Any(), roomID, gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockSlotService(ctrl)
			tt.setup(svc)

			h := NewSlotHandler(svc)

			url := "/rooms/" + tt.roomID + "/slots/list"
			if tt.date != "" {
				url += "?date=" + tt.date
			}
			r := httptest.NewRequest(http.MethodGet, url, nil)
			r = withChiParam(r, roomIDURLParam, tt.roomID)
			w := httptest.NewRecorder()

			h.GetSlots(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
