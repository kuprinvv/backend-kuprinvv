package user

import (
	"test-backend-1-kuprinvv/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.UserRepository = (*repo)(nil)

const (
	tableName = "users"

	idColumn        = "id"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	CreatedAtColumn = "created_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}
