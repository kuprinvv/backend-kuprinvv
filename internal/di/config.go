package di

import (
	"log"
	"test-backend-1-kuprinvv/internal/config"
	"test-backend-1-kuprinvv/internal/config/db"
	"test-backend-1-kuprinvv/internal/config/jwt"
)

type configProvider struct {
	jwtCfg config.JWTConfig
	dbCfg  config.DBConfig
}

func (c *Container) JWTConfig() config.JWTConfig {
	if c.configs.jwtCfg == nil {
		cfg, err := jwt.NewJwtConfig()
		if err != nil {
			log.Fatal(err)
		}
		c.configs.jwtCfg = cfg
	}

	return c.configs.jwtCfg
}

func (c *Container) DBConfig() config.DBConfig {
	if c.configs.dbCfg == nil {
		cfg, err := db.NewDBConfig()
		if err != nil {
			log.Fatal(err)
		}
		c.configs.dbCfg = cfg
	}

	return c.configs.dbCfg
}
