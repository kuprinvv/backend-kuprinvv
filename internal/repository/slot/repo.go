package slot

import (
	"test-backend-1-kuprinvv/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.SlotRepository = (*repo)(nil)

const (
	TableName = "slots"

	IdColumn        = "id"
	RoomIDColumn    = "room_id"
	StartTimeColumn = "start_time"
	EndTimeColumn   = "end_time"
)

type repo struct {
	db *pgxpool.Pool
}

func NewSlotRepository(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}
