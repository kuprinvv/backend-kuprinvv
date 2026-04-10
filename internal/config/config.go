package config

type DBConfig interface {
	DSN() string
}

type JWTConfig interface {
	Token() string
}
