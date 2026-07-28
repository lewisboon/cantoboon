// Package db wires up the Postgres connection pool.
//
// DATABASE_URL examples:
//
//	local dev:  postgres://cantoboon:cantoboon@localhost:5432/cantoboon
//	Cloud Run:  postgres://user:pass@/dbname?host=/cloudsql/PROJECT:REGION:INSTANCE
package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}
	return pgxpool.New(ctx, dsn)
}
