package slot

import (
	"context"
	"log"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r repo) CreateSlots(ctx context.Context, slots []model.Slot) error {
	batch := &pgx.Batch{}

	for _, slot := range slots {
		query := squirrel.Insert(TableName).
			Columns(IdColumn, RoomIDColumn, StartTimeColumn, EndTimeColumn).
			Values(slot.ID, slot.RoomID, slot.StartTime, slot.EndTime).
			Suffix("ON CONFLICT DO NOTHING").
			PlaceholderFormat(squirrel.Dollar)

		sql, args, err := query.ToSql()
		if err != nil {
			log.Printf("failed to generate sql query: %v", err)
		}

		batch.Queue(sql, args...)
	}

	br := r.db.SendBatch(ctx, batch)
	defer func() {
		if err := br.Close(); err != nil {
			log.Printf("batch close failed: %v", err)
		}
	}()

	for range slots {
		_, err := br.Exec()
		if err != nil {
			log.Printf("failed to execute query: %v", err)
			return err
		}
	}

	return nil
}
