package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Entry struct {
	EntryID     int
	Traditional string
	Simplified  string
	Pinyin      string
	Jyutping    string
	Frequency   int
}

// LookupExact returns entries whose traditional or simplified form matches
// query exactly, highest-frequency first.
func LookupExact(ctx context.Context, pool *pgxpool.Pool, query string) ([]Entry, error) {
	rows, err := pool.Query(ctx, `
		SELECT entry_id, traditional, simplified, pinyin, jyutping, frequency
		FROM dictionary.entries
		WHERE traditional = $1 OR simplified = $1
		ORDER BY frequency DESC
	`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.EntryID, &e.Traditional, &e.Simplified, &e.Pinyin, &e.Jyutping, &e.Frequency); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
