-- +goose Up
CREATE TABLE schedule_days (
    id BIGSERIAL PRIMARY KEY,
    schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7)
);

CREATE UNIQUE INDEX idx_unique_schedule_day
    ON schedule_days (schedule_id, day_of_week);
-- +goose Down
DROP TABLE schedule_days;
