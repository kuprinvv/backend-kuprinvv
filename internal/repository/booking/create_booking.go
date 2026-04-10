package booking

import (
	"context"
	"log"
	"strings"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
)

func (r *repo) CreateBooking(ctx context.Context, booking model.Booking) (*model.Booking, error) {
	query := squirrel.Insert(TableName).
		Columns(IdColumn, UserIDColumn, SlotIDColumn, ConferenceLinkColumn).
		Values(booking.ID, booking.UserID, booking.SlotID, booking.ConferenceLink).
		Suffix("RETURNING " + strings.Join([]string{
			IdColumn,
			UserIDColumn,
			SlotIDColumn,
			StatusColumn,
			ConferenceLinkColumn,
			CreatedAtColumn,
		}, ", ")).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to generate query: %v", err)
		return nil, err
	}

	var result model.Booking
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&result.ID,
		&result.UserID,
		&result.SlotID,
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
