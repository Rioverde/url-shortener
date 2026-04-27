-- Create the urls table that maps short codes to original URLs.
CREATE TABLE IF NOT EXISTS urls (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    original_url TEXT     NOT NULL,
    short_code   TEXT     NOT NULL UNIQUE,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Explicit index on short_code for fast lookups during resolution.
-- The UNIQUE constraint above already creates an implicit index, but naming
-- the index explicitly makes intent clear and gives us a stable name to drop.
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
