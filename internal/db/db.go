package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(url string) (*pgxpool.Pool, error) {

	// conn := os.Getenv("postgres://postgres:postgres@localhost:5432/autoflow")

	return pgxpool.New(context.Background(), url)
}
