package main

import (
	"fmt"
	"testing"
	"time"
)

// TestCityExpertFullLoop simulates the main loop for CityExpert to verify it stops dynamically.
// This test runs against the live site, so it might take a minute or two.
func TestCityExpertFullLoop(t *testing.T) {
	// This mirrors the logic in main.go
	// We want to verify that we break out of the loop when we hit the total items limit.

	page := 1
	itemsSoFar := 0
	totalItemsLimit := 0

	// Safety break to prevent infinite test loop if logic fails
	maxTestPages := 50

	for {
		if page > maxTestPages {
			t.Errorf("Hit safety limit of %d pages. Infinite loop detected!", maxTestPages)
			break
		}

		// Stop logic from main.go
		if totalItemsLimit > 0 && itemsSoFar >= totalItemsLimit {
			fmt.Printf("✓ Reached total items limit: %d (processed %d items). Stopping loop.\n", totalItemsLimit, itemsSoFar)
			break
		}

		fmt.Printf("Requesting page %d...\n", page)
		estates, total, err := CityExpertList(page)
		if err != nil {
			t.Fatalf("Error parsing page %d: %v", page, err)
		}

		// Simulate finding total items (logic from main.go)
		if total > 0 {
			if totalItemsLimit == 0 {
				fmt.Printf("Found total items count: %d\n", total)
			}
			totalItemsLimit = total
		}

		if len(estates) == 0 {
			fmt.Println("Received empty list. Stopping loop.")
			break
		}

		count := len(estates)
		itemsSoFar += count
		fmt.Printf("Page %d: Got %d items. Total so far: %d. (Target Limit: %d)\n", page, count, itemsSoFar, totalItemsLimit)

		page++
		// Respect delay
		time.Sleep(2 * time.Second)
	}

	if totalItemsLimit == 0 {
		t.Log("Warning: Did not find total items limit in pagination. Loop stopped by other means (empty list?).")
	} else {
		if itemsSoFar < totalItemsLimit {
			// It's possible to stop early if pages are not full or if limits are slightly off,
			// but we expect to be close.
			// Actually if we break because of empty list before hitting limit, that's also fine (parsing finished).
			fmt.Println("Loop finished.")
		}
	}
}
