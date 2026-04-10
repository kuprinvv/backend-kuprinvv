-- +goose Up
-- +goose StatementBegin
CREATE TABLE schedules (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL UNIQUE,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE schedules;
-- +goose StatementEnd
