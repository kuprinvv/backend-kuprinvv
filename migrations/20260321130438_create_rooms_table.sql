-- +goose Up
-- +goose StatementBegin
CREATE TABLE rooms (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    capacity INT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rooms;
-- +goose StatementEnd
