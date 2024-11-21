-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS library.song (
    song_id serial PRIMARY KEY,
    "group" VARCHAR(50) NOT NULL,
    song VARCHAR(50) NOT NULL,
    release_date date,
    lyrics text,
    link VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS library.song CASCADE;
-- +goose StatementEnd
