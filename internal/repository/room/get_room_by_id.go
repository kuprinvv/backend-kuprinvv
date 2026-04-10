package room

import (
	"context"
	"errors"
	"log"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *repo) GetRoomByID(ctx context.Context, id uuid.UUID) (*model.Room, error) {
	query := squirrel.Select(IdColumn, NameColumn, DescriptionColumn, CapacityColumn, CreatedAtColumn).
		From(TableName).
		Where(squirrel.Eq{IdColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to build sql query: %v", err)
		return nil, err
	}

	var room model.Room

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&room.ID,
		&room.Name,
		&room.Description,
		&room.Capacity,
		&room.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrRoomNotFound
		}
		log.Printf("failed to scan row: %v", err)
		return nil, err
	}

	return &room, nil
}
