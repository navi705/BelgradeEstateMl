package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,sr;q=0.8,rs;q=0.7")
}

func testConnection(t *testing.T, url string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	setHeaders(req)
	req.Header.Set("Referer", "https://www.google.com/")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("status code for %s is %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if len(body) == 0 {
		t.Error("body is empty")
	}
}

func TestFourZidaConnection(t *testing.T) {
	testConnection(t, "https://www.4zida.rs/prodaja-stanova/beograd")
}

func TestFourZidaListFirstPage(t *testing.T) {
	list, _, err := FourZidaList(1)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 0 {
		t.Error("list is empty")
	}
}

func TestFourZidaListSecondPage(t *testing.T) {
	list, _, err := FourZidaList(2)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 0 {
		t.Error("list is empty")
	}
}

func TestFourZidaFloor(t *testing.T) {
	list, _, err := parseWebSiteData("4zida.rs", 1, "https://www.4zida.rs/prodaja-stanova/beograd?sprat_od=-4&sprat_do=-1", "https://www.4zida.rs/prodaja-stanova/beograd?sprat_od=-4&sprat_do=-1&strana=%d", "[test-data='ad-search-card']", parse4ZidaCard, nil)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 0 {
		t.Error("list is empty")
	}
	var countHaveFloor int
	for _, estate := range list {
		if estate.Floor != -1000 {
			countHaveFloor++
		}
		fmt.Println(estate.Floor)
	}
	if countHaveFloor < len(list)/2 {
		t.Error(" less than half estates have floor")
	}

}

func TestFourZidaGenral(t *testing.T) {
	var list []RealEstate

	for i := 1; i < 10; i++ {
		listCommon, _, err := FourZidaList(i)
		if err != nil {
			t.Error(err)
		}
		list = append(list, listCommon...)

		listFloor, _, err := parseWebSiteData("4zida.rs", 1, "https://www.4zida.rs/prodaja-stanova/beograd?sprat_od=-4&sprat_do=-2", "https://www.4zida.rs/prodaja-stanova/beograd?sprat_od=-4&sprat_do=-1&strana=%d", "[test-data='ad-search-card']", parse4ZidaCard, nil)
		if err != nil {
			t.Error(err)
		}
		list = append(list, listFloor...)

		listWhoCreated, _, err := parseWebSiteData("4zida.rs", 1, "https://www.4zida.rs/prodaja-stanova/beograd/investitor?oglasivac=vlasnik", "https://www.4zida.rs/prodaja-stanova/beograd/investitor?oglasivac=vlasnik&strana=%d", "[test-data='ad-search-card']", parse4ZidaCard, nil)
		if err != nil {
			t.Error(err)
		}
		list = append(list, listWhoCreated...)

	}

	if len(list) == 0 {
		t.Error("list is empty")
	}

	var countHaveFloor int
	var countHaveWhoCreated int
	for _, estate := range list {

		if estate.Floor != -1000 {
			countHaveFloor++
		}

		if estate.WhoCreated != 0 {
			countHaveWhoCreated++
		}

		if estate.Price <= 0 {
			t.Error("estate price is less than or equal to 0")
		}

		if estate.SquareMeter <= 0 {
			t.Error("estate square meter is less than or equal to 0")
		}

		if estate.PricePerSquareMeter <= 0 {
			t.Error("estate price per square meter is less than or equal to 0")
		}

		if estate.ParsingDate.IsZero() {
			t.Error("estate parsing date is zero")
		}
		if estate.Source == "" {
			t.Error("estate source is empty")
		}

		if estate.Link == "" {
			t.Error("estate link is empty")
		}

		if estate.FullLocation == "" {
			t.Error("estate full location is empty")
		}
	}
	if countHaveFloor < len(list)/2 {
		t.Error(" less than half estates have floor")
	}
	if countHaveWhoCreated < len(list)/2 {
		t.Error(" less than half estates have who created")
	}
}

func TestHaloOglasiConnection(t *testing.T) {
	testConnection(t, "https://www.halooglasi.com/nekretnine/prodaja-stanova/beograd")
}

func TestHaloOglasiListFirstPage(t *testing.T) {
	list, _, err := HaloOglasiList(1)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 0 {
		t.Error("list is empty")
	}
}

func TestHaloOglasiListSecondPage(t *testing.T) {
	list, _, err := HaloOglasiList(2)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 0 {
		t.Error("list is empty")
	}
}

func TestCityExpertConnection(t *testing.T) {
	testConnection(t, "https://cityexpert.rs/prodaja-nekretnina/beograd?ptId=1")
}

func TestCityExpertListFirstPage(t *testing.T) {
	list, _, err := CityExpertList(1)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 0 {
		t.Error("list is empty")
	}
}

func TestCityExpertListSecondPage(t *testing.T) {
	list, _, err := CityExpertList(2)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 0 {
		t.Error("list is empty")
	}
}

func TestNekretninersConnection(t *testing.T) {
	testConnection(t, "https://www.nekretnine.rs/stambeni-objekti/stanovi/izdavanje-prodaja/prodaja/grad/beograd/lista/")
}

func TestNekretninersListFirstPage(t *testing.T) {
	list, _, err := NekretnineList(1)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 0 {
		t.Error("list is empty")
	}
}

func TestNekretninersListSecondPage(t *testing.T) {
	list, _, err := NekretnineList(2)
	if err != nil {
		t.Error(err)
	}
	if len(list) == 0 {
		t.Error("list is empty")
	}
}

func TestParseNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected int32
	}{
		{"2.397 €/m2", 2397},
		{"57 m2", 57},
		{"1.500 EUR", 1500},
		{"100,50 m2", 101},
		{"  10   ", 10},
		{"No number", 0},
		{"3.630 €/m²", 3630},
		{"1 200", 1200},
		{"1.234,56", 1235}, // Rounding 1234.56 -> 1235
	}

	for _, test := range tests {
		result := parseNumeric(test.input)
		if result != test.expected {
			t.Errorf("parseNumeric(%q) = %d; want %d", test.input, result, test.expected)
		}
	}
}

func TestParseFloorSpecific(t *testing.T) {
	cases := []struct {
		input         string
		expectedFloor float32
		expectedTotal float32
	}{
		{"IV/8", 4, 8},
		{"IV/8.", 4, 8},
		{"IV/8,", 4, 8},
		{"IV/VIII", 4, 8},
	}

	for _, c := range cases {
		floor, total := parseFloor(c.input)
		fmt.Printf("Input: %q, Floor: %f, Total: %f\n", c.input, floor, total)

		if floor != c.expectedFloor {
			t.Errorf("For %q expected Floor %f, got %f", c.input, c.expectedFloor, floor)
		}
		if total != c.expectedTotal {
			t.Errorf("For %q expected Total %f, got %f", c.input, c.expectedTotal, total)
		}
	}
}

func TestParseSerbianRooms(t *testing.T) {
	cases := []struct {
		input    string
		expected float32
	}{
		{"Trosoban stan", 3.0},
		{"Dvosoban", 2.0},
		{"Jednosoban", 1.0},
		{"Garsonjera", 0.5},
		{"Četvorosoban", 4.0},
		{"Petosoban", 5.0},
		{"Petiposoban", 5.5},
		{"  Trosoban  ", 3.0},
		{"Prodaja | Trosoban stan", 3.0},
		{"Unknown", 0},
	}

	for _, c := range cases {
		got := parseSerbianRooms(c.input)
		if got != c.expected {
			t.Errorf("parseSerbianRooms(%q) = %f; want %f", c.input, got, c.expected)
		}
	}
}
