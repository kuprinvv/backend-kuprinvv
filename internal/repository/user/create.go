package user

import (
	"context"
	"log"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
)

func (r *repo) Create(ctx context.Context, user model.User) error {
	query := squirrel.Insert(tableName).
		Columns(idColumn, emailColumn, passwordColumn, roleColumn).
		Values(user.ID, user.Email, user.Password, user.Role).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to build insert query: %v", err)
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		log.Printf("failed to execute insert query: %v", err)
		return err
	}

	return nil
}
