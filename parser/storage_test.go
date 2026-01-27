package main

import (
	"log/slog"
	"testing"
	"time"
)

func TestStorageIntegration(t *testing.T) {
	connStr := "postgres://user:password@localhost:5432/estate_db?sslmode=disable"

	storage, err := NewStorage(connStr)
	if err != nil {
		t.Skipf("Skipping integration test: failed to connect to DB: %v", err)
	}
	defer storage.db.Close()

	// Clean up before test
	_, err = storage.db.Exec("DROP TABLE IF EXISTS estates")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}

	err = storage.Migrate()
	if err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	testEstate := RealEstate{
		Price:               150000,
		Currency:            "EUR",
		PricePerSquareMeter: 2000,
		SquareMeter:         75,
		City:                "Beograd",
		District:            "Novi Beograd",
		Municipality:        "Novi Beograd",
		Street:              "Omladinskih brigada",
		FullLocation:        "Omladinskih brigada, Novi Beograd, Beograd",
		WhoCreated:          Agent,
		QuantityRoom:        3.0,
		Floor:               5,
		FloorTotal:          10,
		Link:                "https://test.com/estate/1",
		ParsingDate:         time.Now(),
		Source:              "test",
	}

	err = storage.SaveEstate(testEstate)
	if err != nil {
		t.Fatalf("SaveEstate failed: %v", err)
	}

	var count int
	err = storage.db.QueryRow("SELECT COUNT(*) FROM estates WHERE link = $1", testEstate.Link).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 estate record, got %d", count)
	}

	testEstate.Price = 160000
	err = storage.SaveEstate(testEstate)
	if err != nil {
		t.Fatalf("SaveEstate update failed: %v", err)
	}

	var newPrice int
	err = storage.db.QueryRow("SELECT price FROM estates WHERE link = $1", testEstate.Link).Scan(&newPrice)
	if err != nil {
		t.Fatalf("Failed to query price: %v", err)
	}

	if newPrice != 160000 {
		t.Errorf("Expected updated price 160000, got %d", newPrice)
	}

	slog.Info("Integration test passed successfully")
}
