package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLStore struct {
	pool *pgxpool.Pool
}

func NewURLStore(databaseURL string) (*URLStore, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing database URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Create table if it doesn't exist
	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS urls (
			short_code VARCHAR(10) PRIMARY KEY,
			long_url TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return &URLStore{pool: pool}, nil
}

func (s *URLStore) Close() {
	s.pool.Close()
}

func (s *URLStore) Set(shortCode, longURL string) error {
	_, err := s.pool.Exec(context.Background(),
		"INSERT INTO urls (short_code, long_url) VALUES ($1, $2)",
		shortCode, longURL,
	)
	if err != nil {
		return fmt.Errorf("error inserting URL: %v", err)
	}
	return nil
}

func (s *URLStore) Get(shortCode string) (string, bool) {
	var longURL string
	err := s.pool.QueryRow(context.Background(),
		"SELECT long_url FROM urls WHERE short_code = $1",
		shortCode,
	).Scan(&longURL)

	if err == pgx.ErrNoRows {
		return "", false
	}
	if err != nil {
		return "", false
	}

	return longURL, true
}

func (s *URLStore) FindByURL(longURL string) (string, bool) {
	var shortCode string
	err := s.pool.QueryRow(context.Background(),
		"SELECT short_code FROM urls WHERE long_url = $1",
		longURL,
	).Scan(&shortCode)

	if err == pgx.ErrNoRows {
		return "", false
	}
	if err != nil {
		return "", false
	}

	return shortCode, true
}
