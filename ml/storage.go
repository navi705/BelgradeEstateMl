package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type WhoCreated int

const (
	Unknown = iota
	Agent
	User
	Investor
)

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format("2006-01-02") + `"`), nil
}

func (d DateOnly) Format(layout string) string {
	return time.Time(d).Format(layout)
}

func (d DateOnly) IsZero() bool {
	return time.Time(d).IsZero()
}

type RealEstate struct {
	Price               int32      `json:"price"`
	Currency            string     `json:"currency"`
	PricePerSquareMeter int32      `json:"price_per_sqm"`
	SquareMeter         int32      `json:"square_meter"`
	City                string     `json:"city"`
	District            string     `json:"district"`
	Municipality        string     `json:"municipality"`
	Street              string     `json:"street"`
	FullLocation        string     `json:"full_location"`
	WhoCreated          WhoCreated `json:"who_created"`
	QuantityRoom        float32    `json:"quantity_room"`
	Floor               float32    `json:"floor"`
	FloorTotal          float32    `json:"floor_total"`
	FloorLabel          string     `json:"floor_label"`
	Link                string     `json:"link"`
	ParsingDate         DateOnly   `json:"parsing_date"`
	Source              string     `json:"source"`
}

func NewConnection(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{db: db}, nil
}

func GetRealEstateWithoutDuplicate(s *Storage, from, to time.Time) ([]RealEstate, error) {
	query := `
	SELECT DISTINCT
	price,
	price_per_sqm,
	square_meter,
	quantity_room,
	FLOOR,
	floor_total,
	district,
	parsing_date
	FROM estates
	WHERE price > 30000 AND currency = 'EUR'
	AND district != '' AND LOWER(district) != 'beograd'
	AND square_meter > 5
	AND price_per_sqm > 300 AND price_per_sqm < 15000
	`

	var args []interface{}
	if !from.IsZero() {
		query += " AND parsing_date >= $1"
		args = append(args, from)
	}
	if !to.IsZero() {
		placeholder := fmt.Sprintf("$%d", len(args)+1)
		query += " AND parsing_date <= " + placeholder
		args = append(args, to)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	defer rows.Close()

	var estates []RealEstate
	for rows.Next() {
		var e RealEstate
		err := rows.Scan(
			&e.Price,
			&e.PricePerSquareMeter,
			&e.SquareMeter,
			&e.QuantityRoom,
			&e.Floor,
			&e.FloorTotal,
			&e.District,
			&e.ParsingDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		e.FloorLabel = GetFloorLabel(e.Floor)
		e.Floor = NormalizeFloorValue(e.Floor, e.FloorTotal)
		e.FloorTotal = NormalizeFloorValue(e.FloorTotal, -5)
		e.District = StandardizeDistrict(e.District)

		estates = append(estates, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return estates, nil
}

func GetDateRange(s *Storage) (time.Time, time.Time, error) {
	var min, max time.Time
	err := s.db.QueryRow("SELECT MIN(parsing_date), MAX(parsing_date) FROM estates").Scan(&min, &max)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to query date range: %w", err)
	}
	return min, max, nil
}
