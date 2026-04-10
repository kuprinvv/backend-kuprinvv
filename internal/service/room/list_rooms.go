package room

import (
	"context"
	"test-backend-1-kuprinvv/internal/model"
)

func (s *serv) ListRooms(ctx context.Context) ([]model.Room, error) {
	return s.roomRepo.ListRooms(ctx)
}
