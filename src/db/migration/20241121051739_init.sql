-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS library
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS library
-- +goose StatementEnd
