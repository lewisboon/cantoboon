-- Run once, after pgloader.load has populated the `dictionary` schema.
--
-- Column names below (pinyin, jyutping, simplified, traditional,
-- definition_text) are taken from the schema shape described while planning
-- this rewrite. Verify them against the real dictionaries.sql after running
-- pgloader — if a name differs (e.g. definitions' text column isn't actually
-- called definition_text), adjust the generated-column expressions below
-- before running this file.

-- jyutping/pinyin syllables aren't English words, so use the 'simple' text
-- search config (tokenize only, no stemming).
ALTER TABLE dictionary.entries
    ADD COLUMN search_vec tsvector
    GENERATED ALWAYS AS (
        to_tsvector('simple', coalesce(pinyin, '') || ' ' || coalesce(jyutping, ''))
    ) STORED;

CREATE INDEX idx_entries_search_vec ON dictionary.entries USING GIN (search_vec);

-- `simplified` lookups would otherwise full-scan; `traditional` should
-- already carry a unique index/constraint carried over from the source
-- SQLite schema via pgloader — confirm with \d dictionary.entries.
CREATE INDEX IF NOT EXISTS idx_entries_simplified ON dictionary.entries (simplified);

-- English meaning search benefits from stemming.
ALTER TABLE dictionary.definitions
    ADD COLUMN search_vec tsvector
    GENERATED ALWAYS AS (
        to_tsvector('english', coalesce(definition_text, ''))
    ) STORED;

CREATE INDEX idx_definitions_search_vec ON dictionary.definitions USING GIN (search_vec);
