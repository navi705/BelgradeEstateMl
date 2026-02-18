package main

import "testing"

func TestStandardizeDistrict(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Mirijevo I", "Zvezdara"},
		{"Mirijevo II", "Zvezdara"},
		{"Kalenić pijaca", "Vračar"},
		{"Blok 65", "Novi Beograd"},
		{"Novi Beograd Blok 65", "Novi Beograd"},
		{"Bežanijska kosa 1", "Novi Beograd"},
		{"Dedinje", "Savski venac"},
		{"Beograd na vodi", "Savski venac"},
		{"Karaburma", "Palilula"},
		{"Medaković III", "Voždovac"},
		{"Braće Jerković", "Voždovac"},
		{"Altina", "Zemun"},
		{"Batajnica", "Zemun"},
		{"Dorćol", "Stari Grad"},
		{"Banovo brdo", "Čukarica"},
		{"Miljakovac I", "Rakovica"},
		{"Nepoznat Kraj", "Nepoznat Kraj"},
		{"", "Unknown"},
		{"  Mirijevo I  ", "Zvezdara"},
	}

	for _, tt := range tests {
		got := StandardizeDistrict(tt.input)
		if got != tt.expected {
			t.Errorf("StandardizeDistrict(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestGetAllStandardizedDistricts(t *testing.T) {
	districts := GetAllStandardizedDistricts()
	if len(districts) == 0 {
		t.Error("Expected non-empty list of districts")
	}

	found := false
	for _, d := range districts {
		if d == "Novi Beograd" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected 'Novi Beograd' in the list of districts")
	}
}
