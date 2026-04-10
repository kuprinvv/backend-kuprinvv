package room

import "test-backend-1-kuprinvv/internal/service"

type Handler struct {
	roomService service.RoomService
}

func NewRoomHandler(roomService service.RoomService) *Handler {
	return &Handler{
		roomService: roomService,
	}
}
