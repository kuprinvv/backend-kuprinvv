-- +goose Up
-- +goose StatementBegin
CREATE TABLE slots (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,

    CONSTRAINT valid_slot CHECK (start_time < end_time)
);

CREATE UNIQUE INDEX idx_slots_unique
    ON slots (room_id, start_time);

CREATE INDEX idx_slots_room_time
    ON slots (room_id, start_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE slots;
-- +goose StatementEnd
