package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func getConnection() (*Storage, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Warning: error loading .env file: %v\n", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	fmt.Println("Connecting to:", dbURL)
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	return NewConnection(dbURL)
}

func TestGetRealEstateWithoutDuplicate(t *testing.T) {
	storage, err := getConnection()
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	result, err := GetRealEstateWithoutDuplicate(storage, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("failed to get real estate without duplicate: %v", err)
	}
	defer storage.db.Close()

	if len(result) == 0 {
		t.Fatal("expected at least one real estate without duplicate")
	}
	fmt.Println("Number of real estates without duplicate:", len(result))

	correlationTable(result, "All Districts")

	district := "Centar"
	vracarEstates := FilterByDistrict(result, district)
	fmt.Printf("\n--- Correlation for %s (%d estates) ---\n", district, len(vracarEstates))
	correlationTable(vracarEstates, district)
}

func correlationTable(estates []RealEstate, title string) {

	matrix := CorrelationMatrix(estates)
	labels := []string{"Price", "Sqm", "Rooms", "Floor", "FloorTotal"}

	fmt.Printf("%-12s", "")
	for _, l := range labels {
		fmt.Printf("%-10s", l)
	}
	fmt.Println()

	for i, row := range matrix {
		fmt.Printf("%-12s", labels[i])
		for _, val := range row {
			fmt.Printf("%-10.4f", val)
		}
		fmt.Println()
	}
}
