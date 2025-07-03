-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS list_entries (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    order_index INTEGER NOT NULL,
    list_id BIGSERIAL NOT NULL REFERENCES lists(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE list_entries;
-- +goose StatementEnd