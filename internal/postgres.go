package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDb(ctx context.Context, user, password, host, database string) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, password, host, database)
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
