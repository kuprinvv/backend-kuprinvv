package booking

import (
	"test-backend-1-kuprinvv/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.BookingRepository = (*repo)(nil)

const (
	TableName = "bookings"

	IdColumn             = "id"
	UserIDColumn         = "user_id"
	SlotIDColumn         = "slot_id"
	StatusColumn         = "status"
	ConferenceLinkColumn = "conference_link"
	CreatedAtColumn      = "created_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}
