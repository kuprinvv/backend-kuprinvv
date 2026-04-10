package room

import (
	"context"
	"log"
	"strings"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
)

func (r *repo) CreateRoom(ctx context.Context, room model.Room) (*model.Room, error) {
	query := squirrel.Insert(TableName).
		Columns(IdColumn, NameColumn, DescriptionColumn, CapacityColumn).
		Values(room.ID, room.Name, room.Description, room.Capacity).
		Suffix("RETURNING " + strings.Join([]string{
			IdColumn,
			NameColumn,
			DescriptionColumn,
			CapacityColumn,
			CreatedAtColumn,
		}, ", ")).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to build sql query: %v", err)
		return nil, err
	}

	var result model.Room
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.Capacity,
		&result.CreatedAt,
	)
	if err != nil {
		log.Printf("failed to scan row: %v", err)
		return nil, err
	}

	return &result, nil
}
