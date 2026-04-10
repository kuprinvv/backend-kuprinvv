package service

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIsUniqueViolation(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "уникальное нарушение (23505)",
			err:  &pgconn.PgError{Code: "23505"},
			want: true,
		},
		{
			name: "другая ошибка pgx",
			err:  &pgconn.PgError{Code: "42P01"},
			want: false,
		},
		{
			name: "обычная ошибка",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsUniqueViolation(tt.err))
		})
	}
}
