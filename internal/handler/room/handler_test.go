package room

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"test-backend-1-kuprinvv/internal/mocks"
	"test-backend-1-kuprinvv/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateRoom(t *testing.T) {
	roomID := uuid.New()
	name := "Переговорка 1"

	tests := []struct {
		name       string
		body       any
		setup      func(svc *mocks.MockRoomService)
		wantStatus int
	}{
		{
			name: "успешное создание",
			body: map[string]any{"name": name},
			setup: func(svc *mocks.MockRoomService) {
				svc.EXPECT().
					CreateRoom(gomock.Any(), gomock.Any()).
					Return(&model.Room{ID: roomID, Name: name}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "невалидный JSON",
			body:       "not json",
			setup:      func(svc *mocks.MockRoomService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "пустое имя",
			body:       map[string]any{"name": ""},
			setup:      func(svc *mocks.MockRoomService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "внутренняя ошибка",
			body: map[string]any{"name": name},
			setup: func(svc *mocks.MockRoomService) {
				svc.EXPECT().
					CreateRoom(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockRoomService(ctrl)
			tt.setup(svc)

			h := NewRoomHandler(svc)

			rawBody, _ := json.Marshal(tt.body)
			r := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewReader(rawBody))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.CreateRoom(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestListRooms(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(svc *mocks.MockRoomService)
		wantStatus int
	}{
		{
			name: "успешно",
			setup: func(svc *mocks.MockRoomService) {
				svc.EXPECT().
					ListRooms(gomock.Any()).
					Return([]model.Room{{ID: uuid.New(), Name: "Room 1"}}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "пустой список",
			setup: func(svc *mocks.MockRoomService) {
				svc.EXPECT().ListRooms(gomock.Any()).Return([]model.Room{}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "внутренняя ошибка",
			setup: func(svc *mocks.MockRoomService) {
				svc.EXPECT().ListRooms(gomock.Any()).Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockRoomService(ctrl)
			tt.setup(svc)

			h := NewRoomHandler(svc)
			r := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
			w := httptest.NewRecorder()

			h.ListRooms(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
