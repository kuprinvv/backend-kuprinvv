package user

import (
	"context"
	"log"
	"test-backend-1-kuprinvv/internal/model"

	"github.com/Masterminds/squirrel"
)

func (r *repo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := squirrel.Select(idColumn, emailColumn, passwordColumn, roleColumn, CreatedAtColumn).
		From(tableName).
		Where(squirrel.Eq{emailColumn: email}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Printf("failed to build sql query: %v", err)
		return nil, err
	}

	row := r.db.QueryRow(ctx, sql, args...)

	var u model.User
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.CreatedAt); err != nil {
		log.Printf("failed to scan row: %v", err)
		return nil, err
	}

	return &u, nil
}
