package internal

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDb(url string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return db, err
}
