-- +goose Up
-- +goose StatementBegin
CREATE TABLE bookings (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('active', 'cancelled')) DEFAULT 'active',
    conference_link TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_unique_active_booking
    ON bookings (slot_id)
    WHERE status = 'active';

CREATE INDEX idx_bookings_user
    ON bookings (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bookings;
-- +goose StatementEnd
