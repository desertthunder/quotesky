-- LibSQL Schema
--
CREATE TABLE apps IF NOT EXISTS (
    handle TEXT PRIMARY KEY NOT NULL,
    token TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
