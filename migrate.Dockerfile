FROM golang:1.26

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app
COPY migrations ./migrations

CMD goose -dir migrations postgres "$DB_DSN" up -v && exit 0