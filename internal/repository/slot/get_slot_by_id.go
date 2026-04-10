package slot

import (
	"context"
	"errors"
	"log"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r repo) GetSlotByID(ctx context.Context, id uuid.UUID) (*model.Slot, error) {
	query := squirrel.Select(IdColumn, RoomIDColumn, StartTimeColumn, EndTimeColumn).
		From(TableName).
		Where(squirrel.Eq{IdColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to build sql query: %v", err)
		return nil, err
	}

	row := r.db.QueryRow(ctx, sql, args...)

	var s model.Slot

	if err = row.Scan(&s.ID, &s.RoomID, &s.StartTime, &s.EndTime); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrSlotNotFound
		}
		log.Printf("failed to scan row: %v", err)
		return nil, err
	}

	return &s, nil
}
