package room

import (
	"context"
	"log"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
)

func (r *repo) ListRooms(ctx context.Context) ([]model.Room, error) {
	query := squirrel.Select(IdColumn, NameColumn, DescriptionColumn, CapacityColumn, CreatedAtColumn).
		From(TableName).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to build sql query: %v", err)
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		log.Printf("failed to query rooms: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rooms []model.Room

	for rows.Next() {
		var room model.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt); err != nil {
			log.Printf("failed to scan row: %v", err)
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows iteration error: %v", err)
		return nil, err
	}

	return rooms, nil
}
