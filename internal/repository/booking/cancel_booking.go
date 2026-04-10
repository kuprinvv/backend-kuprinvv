package booking

import (
	"context"
	"log"
	"strings"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repo) CancelBooking(ctx context.Context, bookingID uuid.UUID) (*model.Booking, error) {
	query := squirrel.Update(TableName).
		Set(StatusColumn, "cancelled").
		Where(squirrel.Eq{IdColumn: bookingID}).
		Suffix("RETURNING " + strings.Join([]string{
			IdColumn,
			SlotIDColumn,
			UserIDColumn,
			StatusColumn,
			ConferenceLinkColumn,
			CreatedAtColumn,
		}, ", ")).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to generate sql: %v", err)
		return nil, err
	}

	var result model.Booking
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&result.ID,
		&result.SlotID,
		&result.UserID,
		&result.Status,
		&result.ConferenceLink,
		&result.CreatedAt,
	)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return nil, err
	}

	return &result, nil
}
