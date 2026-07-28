Migration from the local `dictionaries.db` (SQLite) into Cloud SQL for PostgreSQL.

Run in order, against the real `dictionaries.db` (this repo doesn't contain
that file — it's local-only, see the plan for why):

1. `pgloader pgloader.load` — copies the base tables/data into the
   `dictionary` schema. Edit the `FROM sqlite:///...` path first. Skips the
   `entries_fts`/`definitions_fts` FTS5 virtual tables (not real tables, no
   Postgres equivalent).
2. `psql <connection> -f 002_search_indexes.sql` — adds generated `tsvector`
   columns + GIN indexes (replacing FTS5/bm25) and `idx_entries_simplified`.
   Column names in this file are assumptions based on the schema shape
   described while planning this rewrite — verify against the actual
   `dictionaries.sql` after step 1 and adjust if any name differs.
3. Verify row counts match the source SQLite file, and spot-check a few
   known words (e.g. 廣場, 你好) return sensible results via `search_vec @@
   to_tsquery(...)`.

`testdata/fixture.sql` is synthetic sample data shaped like pgloader's output,
used to smoke-test `002_search_indexes.sql` locally without the real 293MB
file. It is not part of the real migration.
