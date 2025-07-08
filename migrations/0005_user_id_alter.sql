-- +goose Up
-- +goose StatementBegin
ALTER TABLE lists
ADD COLUMN user_id BIGSERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE lists DROP COLUMN user_id;
-- +goose StatementEnd