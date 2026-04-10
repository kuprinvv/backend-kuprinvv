package schedule

import (
	"context"
	"log"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repo) GetByRoom(ctx context.Context, roomID uuid.UUID) (*model.Schedule, error) {
	query := squirrel.Select(
		TableName+"."+IdColumn,
		TableName+"."+RoomIDColumn,
		TableName+"."+StartTimeColumn,
		TableName+"."+EndTimeColumn,
		DaysTable+"."+DayOfWeekColumn,
	).
		From(TableName).
		Join(DaysTable + " ON " + DaysTable + "." + DaysScheduleIdColumn + " = " + TableName + "." + IdColumn).
		Where(squirrel.Eq{TableName + "." + RoomIDColumn: roomID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var result model.Schedule
	dayMap := map[int]bool{}

	for rows.Next() {
		var d int
		if err := rows.Scan(
			&result.ID,
			&result.RoomID,
			&result.StartTime,
			&result.EndTime,
			&d,
		); err != nil {
			log.Printf("failed to scan row: %v", err)
			return nil, err
		}
		if !dayMap[d] {
			result.DaysOfWeek = append(result.DaysOfWeek, d)
			dayMap[d] = true
		}
	}

	if result.ID == uuid.Nil {
		return nil, nil
	}

	return &result, nil
}
