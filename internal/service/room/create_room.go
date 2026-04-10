package room

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/google/uuid"
)

func (s *serv) CreateRoom(ctx context.Context, room model.Room) (*model.Room, error) {
	room.ID = uuid.New()

	return s.roomRepo.CreateRoom(ctx, room)
}
