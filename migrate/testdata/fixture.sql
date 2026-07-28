-- Synthetic sample data shaped like what pgloader.load produces from the
-- real dictionaries.db. Used only to smoke-test 002_search_indexes.sql
-- locally — NOT a substitute for running the real migration.

CREATE SCHEMA IF NOT EXISTS dictionary;

CREATE TABLE dictionary.entries (
    entry_id    SERIAL PRIMARY KEY,
    traditional TEXT UNIQUE NOT NULL,
    simplified  TEXT,
    pinyin      TEXT,
    jyutping    TEXT,
    frequency   INTEGER DEFAULT 0
);

CREATE TABLE dictionary.definitions (
    definition_id   SERIAL PRIMARY KEY,
    entry_id        INTEGER REFERENCES dictionary.entries(entry_id),
    definition_text TEXT
);

INSERT INTO dictionary.entries (traditional, simplified, pinyin, jyutping, frequency) VALUES
    ('廣場', '广场', 'guang3 chang3', 'gwong2 coeng4', 100),
    ('你好', '你好', 'ni3 hao3',      'nei5 hou2',    500),
    ('廣東話', '广东话', 'guang3 dong1 hua4', 'gwong2 dung1 waa2', 300);

INSERT INTO dictionary.definitions (entry_id, definition_text) VALUES
    (1, 'plaza; public square'),
    (2, 'hello; how are you'),
    (3, 'Cantonese (spoken language)');
