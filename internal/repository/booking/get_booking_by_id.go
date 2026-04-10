package booking

import (
	"context"
	"errors"
	"log"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *repo) GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error) {
	query := squirrel.Select(IdColumn, UserIDColumn, SlotIDColumn, StatusColumn, ConferenceLinkColumn, CreatedAtColumn).
		From(TableName).
		Where(squirrel.Eq{IdColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to generate sql query: %v", err)
		return nil, err
	}

	var b model.Booking
	if err = r.db.QueryRow(ctx, sql, args...).
		Scan(&b.ID, &b.UserID, &b.SlotID, &b.Status, &b.ConferenceLink, &b.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrBookingNotFound
		}
		log.Printf("failed to scan row: %v", err)
		return nil, err
	}

	return &b, nil
}
