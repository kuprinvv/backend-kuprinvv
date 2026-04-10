package booking

import (
	"context"
	"log"
	"test-backend-1-kuprinvv/internal/model"
	"test-backend-1-kuprinvv/internal/repository/slot"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repo) GetFutureBookingsByUser(ctx context.Context, userID uuid.UUID) ([]model.Booking, error) {
	query := squirrel.Select(
		TableName+"."+IdColumn,
		TableName+"."+SlotIDColumn,
		TableName+"."+UserIDColumn,
		TableName+"."+StatusColumn,
		TableName+"."+ConferenceLinkColumn,
		TableName+"."+CreatedAtColumn,
	).
		From(TableName).
		Join(slot.TableName + " ON " + slot.TableName + "." + slot.IdColumn + " = " + TableName + "." + SlotIDColumn).
		Where(squirrel.Eq{TableName + "." + UserIDColumn: userID}).
		Where(slot.TableName + "." + slot.StartTimeColumn + " > NOW()").
		OrderBy(slot.TableName + "." + slot.StartTimeColumn).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to generate sql for GetFutureBookingsByUser: %+v\n", err)
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		log.Printf("failed to query GetFutureBookingsByUser: %+v\n", err)
		return nil, err
	}
	defer rows.Close()

	var result []model.Booking
	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(&b.ID, &b.SlotID, &b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt); err != nil {
			log.Printf("failed to scan GetFutureBookingsByUser: %+v\n", err)
			return nil, err
		}
		result = append(result, b)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows iteration error: %v", err)
		return nil, err
	}

	return result, nil
}
