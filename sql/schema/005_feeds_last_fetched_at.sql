-- +goose UP

ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMPTZ;

-- +goose DOWN

ALTER TABLE feeds DROP COLUMN IF EXISTS last_fetched_at;