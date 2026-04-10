package db

import (
	"errors"
	"os"
	"test-backend-1-kuprinvv/internal/config"
)

var _ config.DBConfig = (*dbConf)(nil)

const (
	dsnEnvName = "DB_DSN"
)

type dbConf struct {
	dsn string
}

func NewDBConfig() (*dbConf, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("database dsn env var not set")
	}

	return &dbConf{dsn: dsn}, nil
}

func (d *dbConf) DSN() string {
	return d.dsn
}
