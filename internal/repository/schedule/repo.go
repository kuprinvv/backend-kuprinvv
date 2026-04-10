package schedule

import (
	"test-backend-1-kuprinvv/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.ScheduleRepository = (*repo)(nil)

const (
	TableName = "schedules"

	IdColumn        = "id"
	RoomIDColumn    = "room_id"
	DayOfWeekColumn = "day_of_week"
	StartTimeColumn = "start_time"
	EndTimeColumn   = "end_time"
	CreatedAtColumn = "created_at"

	DaysTable            = "schedule_days"
	DaysIdColumn         = "id"
	DaysScheduleIdColumn = "schedule_id"
	DaysDayOfWeekColumn  = "day_of_week"
)

type repo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}
