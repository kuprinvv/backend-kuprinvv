package di

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type dbProvider struct {
	pool *pgxpool.Pool
}

func (c *Container) DB(ctx context.Context) *pgxpool.Pool {
	if c.db.pool == nil {
		pool, err := pgxpool.New(ctx, c.DBConfig().DSN())
		if err != nil {
			log.Fatal(err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatal(err)
		}

		c.db.pool = pool
	}

	return c.db.pool
}
