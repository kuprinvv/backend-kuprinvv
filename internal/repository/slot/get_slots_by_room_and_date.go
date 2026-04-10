package slot

import (
	"context"
	"log"
	"test-backend-1-kuprinvv/internal/model"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r repo) GetSlotsByRoomAndDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]model.Slot, error) {
	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	query := squirrel.Select(TableName+"."+IdColumn, TableName+"."+RoomIDColumn, TableName+"."+StartTimeColumn, TableName+"."+EndTimeColumn).
		From(TableName).
		LeftJoin("bookings ON bookings.slot_id = " + TableName + "." + IdColumn + " AND bookings.status = 'active'").
		Where("bookings.slot_id IS NULL").
		Where(squirrel.Eq{TableName + "." + RoomIDColumn: roomID}).
		Where(squirrel.GtOrEq{TableName + "." + StartTimeColumn: start}).
		Where(squirrel.Lt{TableName + "." + StartTimeColumn: end}).
		OrderBy(TableName + "." + StartTimeColumn).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to generate sql: %v", err)
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var slots []model.Slot

	for rows.Next() {
		var s model.Slot
		if err = rows.Scan(&s.ID, &s.RoomID, &s.StartTime, &s.EndTime); err != nil {
			log.Printf("failed to scan row: %v", err)
			return nil, err
		}
		slots = append(slots, s)
	}

	if err = rows.Err(); err != nil {
		log.Printf("rows iteration error: %v", err)
		return nil, err
	}

	return slots, nil
}
