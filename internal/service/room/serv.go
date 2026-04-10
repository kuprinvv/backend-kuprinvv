package room

import (
	"test-backend-1-kuprinvv/internal/repository"
	"test-backend-1-kuprinvv/internal/service"
)

var _ service.RoomService = (*serv)(nil)

type serv struct {
	roomRepo repository.RoomRepository
}

func NewRoomService(roomRepo repository.RoomRepository) *serv {
	return &serv{roomRepo: roomRepo}
}
