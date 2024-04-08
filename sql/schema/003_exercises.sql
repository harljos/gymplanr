-- +goose Up
CREATE TABLE exercises (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    muscle TEXT NOT NULL,
    repetitions INTEGER NOT NULL,
    instructions TEXT NOT NULL,
    day_id UUID NOT NULL REFERENCES days(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE exercises;