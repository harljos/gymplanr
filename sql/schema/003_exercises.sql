-- +goose Up
CREATE TABLE exercises (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    muscle TEXT NOT NULL,
    sets INTEGER,
    repetitions INTEGER,
    exercise_duration INTEGER,
    instructions TEXT NOT NULL,
    exercise_type TEXT NOT NULL,
    day_id UUID NOT NULL REFERENCES days(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE exercises;