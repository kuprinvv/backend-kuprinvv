package booking

import (
	"context"
	"log"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/internal/repository/slot"

	"github.com/Masterminds/squirrel"
)

func (r *repo) ListBookings(ctx context.Context, limit, offset int) ([]model.Booking, int, error) {
	query := squirrel.Select(
		TableName+"."+IdColumn,
		TableName+"."+UserIDColumn,
		TableName+"."+SlotIDColumn,
		TableName+"."+StatusColumn,
		TableName+"."+ConferenceLinkColumn,
		TableName+"."+CreatedAtColumn,
		"COUNT(*) OVER() AS total",
	).
		From(TableName).
		Join(slot.TableName + " ON " + slot.TableName + "." + slot.IdColumn + " = " + TableName + "." + SlotIDColumn).
		OrderBy(slot.TableName + "." + slot.StartTimeColumn + " DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var result []model.Booking
	var total int

	for rows.Next() {
		var b model.Booking
		if err = rows.Scan(
			&b.ID,
			&b.UserID,
			&b.SlotID,
			&b.Status,
			&b.ConferenceLink,
			&b.CreatedAt,
			&total,
		); err != nil {
			log.Printf("failed to scan row: %v", err)
			return nil, 0, err
		}
		result = append(result, b)
	}

	if err = rows.Err(); err != nil {
		log.Printf("rows iteration error: %v", err)
		return nil, 0, err
	}

	return result, total, nil
}
