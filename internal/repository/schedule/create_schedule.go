package schedule

import (
	"context"
	"errors"
	"log"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *repo) CreateSchedule(ctx context.Context, schedule model.Schedule) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("rollback failed: %v", err)
		}
	}()

	scheduleQuery := squirrel.Insert(TableName).
		Columns(IdColumn, RoomIDColumn, StartTimeColumn, EndTimeColumn).
		Values(schedule.ID, schedule.RoomID, schedule.StartTime, schedule.EndTime).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := scheduleQuery.ToSql()
	if err != nil {
		log.Printf("failed to build schedule query: %v", err)
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		log.Printf("failed to insert schedule: %v", err)
		return err
	}

	for _, d := range schedule.DaysOfWeek {
		dayQuery := squirrel.Insert(DaysTable).
			Columns(DaysScheduleIdColumn, DayOfWeekColumn).
			Values(schedule.ID, d).
			PlaceholderFormat(squirrel.Dollar)

		sql, args, err := dayQuery.ToSql()
		if err != nil {
			log.Printf("failed to build day query: %v", err)
			return err
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			log.Printf("failed to insert schedule day: %v", err)
			return err
		}
	}

	return tx.Commit(ctx)
}
