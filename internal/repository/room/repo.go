package room

import (
	"test-backend-1-kuprinvv/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.RoomRepository = (*repo)(nil)

const (
	TableName = "rooms"

	IdColumn          = "id"
	NameColumn        = "name"
	DescriptionColumn = "description"
	CapacityColumn    = "capacity"
	CreatedAtColumn   = "created_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRoomRepository(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}
