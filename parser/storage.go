package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS estates (
		id SERIAL PRIMARY KEY,
		price INTEGER,
		currency TEXT,
		price_per_sqm INTEGER,
		square_meter INTEGER,
		city TEXT,
		district TEXT,
		municipality TEXT,
		street TEXT,
		full_location TEXT,
		who_created INTEGER,
		quantity_room REAL,
		floor REAL,
		floor_total REAL,
		link TEXT UNIQUE,
		parsing_date TIMESTAMP,
		source TEXT
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}
	slog.Info("Database migration completed successfully")
	return nil
}

func (s *Storage) SaveEstate(e RealEstate) error {
	query := `
	INSERT INTO estates (
		price, currency, price_per_sqm, square_meter, city, district, municipality, street, 
		full_location, who_created, quantity_room, floor, floor_total, link, parsing_date, source
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, 
		$9, $10, $11, $12, $13, $14, $15, $16
	) ON CONFLICT (link) DO UPDATE SET
		price = EXCLUDED.price,
		parsing_date = EXCLUDED.parsing_date,
		price_per_sqm = EXCLUDED.price_per_sqm;
	`

	_, err := s.db.Exec(query,
		e.Price, e.Currency, e.PricePerSquareMeter, e.SquareMeter, e.City, e.District, e.Municipality, e.Street,
		e.FullLocation, e.WhoCreated, e.QuantityRoom, e.Floor, e.FloorTotal, e.Link, time.Now(), e.Source,
	)

	if err != nil {
		slog.Error("failed to save estate", "link", e.Link, "error", err)
		return fmt.Errorf("failed to save estate: %w", err)
	}
	slog.Debug("estate saved successfully", "link", e.Link, "source", e.Source)
	return nil
}
