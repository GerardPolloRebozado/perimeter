-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE hooks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    payload     JSONB NOT NULL
);

-- +goose Down
DROP TABLE hooks;
