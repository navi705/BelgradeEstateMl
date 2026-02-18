package main

import (
	"math"
	"testing"
)

func TestAvg(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"empty slice", []float64{}, 0},
		{"single element", []float64{5}, 5},
		{"multiple elements", []float64{1, 2, 3, 4, 5}, 3},
		{"negative elements", []float64{-1, 1}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Avg(tt.input); got != tt.expected {
				t.Errorf("Avg() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMedian(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"even", []float64{1, 2, 3, 4}, 2.5},
		{"odd", []float64{1, 2, 3, 4, 5}, 3},
		{"unsorted", []float64{5, 1, 3}, 3},
		{"empty", []float64{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Median(tt.input); got != tt.expected {
				t.Errorf("Median() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMode(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected []float64
	}{
		{"single mode", []float64{1, 2, 2, 3}, []float64{2}},
		{"bi-modal", []float64{1, 1, 2, 2, 3}, []float64{1, 2}},
		{"no repeats", []float64{1, 2, 3}, []float64{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mode(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("Mode() count = %v, want %v", len(got), len(tt.expected))
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Mode()[%d] = %v, want %v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestPercentile(t *testing.T) {
	data := []float64{15, 20, 35, 40, 50}
	tests := []struct {
		name     string
		p        float64
		expected float64
	}{
		{"0th", 0, 15},
		{"40th", 40, 29},
		{"50th", 50, 35},
		{"100th", 100, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Percentile(data, tt.p); got != tt.expected {
				t.Errorf("Percentile(%v) = %v, want %v", tt.p, got, tt.expected)
			}
		})
	}
}

func TestQuartile(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if got := Quartile(data, 2); got != 5.5 {
		t.Errorf("Quartile(2) = %v, want 5.5", got)
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		y        int
		expected float64
	}{
		{"2^0", 2, 0, 1},
		{"2^3", 2, 3, 8},
		{"-2^2", -2, 2, 4},
		{"-2^3", -2, 3, -8},
		{"0^5", 0, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pow(tt.x, tt.y); got != tt.expected {
				t.Errorf("Pow() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"sqrt 0", 0, 0},
		{"sqrt 1", 1, 1},
		{"sqrt 4", 4, 2},
		{"sqrt 9", 9, 3},
		{"negative input", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sqrt(tt.input)
			if math.Abs(got-tt.expected) > 0.000001 {
				t.Errorf("Sqrt() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		precision int
		expected  float64
	}{
		{"round 1.234 to 2", 1.234, 2, 1.23},
		{"round 1.235 to 2", 1.235, 2, 1.24},
		{"round 1.5 to 0", 1.5, 0, 2},
		{"round 1.4 to 0", 1.4, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Round(tt.input, tt.precision); got != tt.expected {
				t.Errorf("Round() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"positive", 5, 5},
		{"negative", -5, 5},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.input); got != tt.expected {
				t.Errorf("Abs() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCorrleation(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{1, 2, 3, 4, 5}
	got := Correlation(x, y)
	if got == 0 {
		t.Error("Expected non-zero correlation for identical slices")
	}
}

func TestIQR(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := 4.5
	if got := IQR(data); got != expected {
		t.Errorf("IQR() = %v, want %v", got, expected)
	}
}

func TestGetOutlierBounds(t *testing.T) {
	data := []float64{1, 2, 3, 4, 100}
	lower, upper := GetOutlierBounds(data)
	if lower != -1 || upper != 7 {
		t.Errorf("GetOutlierBounds() = (%v, %v), want (-1, 7)", lower, upper)
	}
}

func TestHistogram(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 10}
	bins := Histogram(data, 2, false)

	t.Logf("Histogram result for %v:", data)
	for i, bin := range bins {
		t.Logf("  Bin %d: [%v - %v] Count: %v", i, bin["from"], bin["to"], bin["count"])
	}

	if len(bins) != 2 {
		t.Fatalf("Expected 2 bins, got %d", len(bins))
	}
	if bins[0]["count"].(int) != 5 {
		t.Errorf("Bin 0 count = %d, want 5", bins[0]["count"])
	}
	if bins[1]["count"].(int) != 1 {
		t.Errorf("Bin 1 count = %d, want 1", bins[1]["count"])
	}

	roundedBins := Histogram(data, 2, true)
	if roundedBins[0]["from"].(float64) != 1 {
		t.Errorf("Expected rounded from = 1, got %v", roundedBins[0]["from"])
	}
}

func TestVarianceAndStdDev(t *testing.T) {
	data := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	variance := Variance(data)
	stddev := StdDev(data)

	if math.Abs(variance-4) > 0.0001 {
		t.Errorf("Variance = %v, want 4", variance)
	}
	if math.Abs(stddev-2) > 0.0001 {
		t.Errorf("StdDev = %v, want 2", stddev)
	}
}

func TestIsNormalDistribution(t *testing.T) {
	data := []float64{
		10, 12, 15, 17, 18, 19, 20, 21, 22, 23, 25, 28, 30,
		10, 12, 15, 17, 18, 19, 20, 21, 22, 23, 25, 28, 30,
		20, 20, 20, 20, 20,
	}

	isNormal := IsNormalDistribution(data)
	t.Logf("Is normal: %v", isNormal)
}

func TestZScore(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected []float64
	}{
		{
			"standard dataset",
			[]float64{2, 4, 4, 4, 5, 5, 7, 9},
			[]float64{-1.5, -0.5, -0.5, -0.5, 0, 0, 1, 2},
		},
		{
			"constant dataset",
			[]float64{5, 5, 5},
			[]float64{0, 0, 0},
		},
		{
			"empty slice",
			[]float64{},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ZScore(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("ZScore() length = %v, want %v", len(got), len(tt.expected))
			}
			for i := range got {
				if math.Abs(got[i]-tt.expected[i]) > 0.0001 {
					t.Errorf("ZScore()[%d] = %v, want %v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestLinearRegression(t *testing.T) {
	tests := []struct {
		name          string
		x             []float64
		y             []float64
		wantSlope     float64
		wantIntercept float64
	}{
		{
			"simple linear correlation",
			[]float64{1, 2, 3, 4, 5},
			[]float64{2, 4, 6, 8, 10},
			2.0,
			0.0,
		},
		{
			"with intercept",
			[]float64{1, 2, 3},
			[]float64{4, 6, 8},
			2.0,
			2.0,
		},
		{
			"mismatched length",
			[]float64{1, 2},
			[]float64{1},
			0,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slope, intercept := LinearRegression(tt.x, tt.y)
			if math.Abs(slope-tt.wantSlope) > 0.0001 {
				t.Errorf("LinearRegression() slope = %v, want %v", slope, tt.wantSlope)
			}
			if math.Abs(intercept-tt.wantIntercept) > 0.0001 {
				t.Errorf("LinearRegression() intercept = %v, want %v", intercept, tt.wantIntercept)
			}
		})
	}
}

func TestGetSigmaBounds(t *testing.T) {
	data := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	// mean=5, stddev=2
	min, max := GetSigmaBounds(data, 1)
	if math.Abs(min-3) > 0.0001 || math.Abs(max-7) > 0.0001 {
		t.Errorf("GetSigmaBounds() = (%v, %v), want (3, 7)", min, max)
	}
}

func TestSolveOLS(t *testing.T) {
	X := [][]float64{
		{1, 1},
		{1, 2},
		{1, 3},
	}
	Y := []float64{2, 4, 6}
	weights := SolveOLS(X, Y)
	if math.Abs(weights[0]-0) > 0.0001 || math.Abs(weights[1]-2) > 0.0001 {
		t.Errorf("SolveOLS() = %v, want [0, 2]", weights)
	}
}

func TestRSquared(t *testing.T) {
	actual := []float64{1, 2, 3, 4, 5}
	predicted := []float64{1.1, 1.9, 3.2, 3.8, 5.1}
	r2 := RSquared(actual, predicted)
	if r2 < 0.9 || r2 > 1.0 {
		t.Errorf("RSquared() = %v, want approx 0.99", r2)
	}

	adjR2 := AdjustedRSquared(r2, 5, 1)
	if adjR2 >= r2 {
		t.Errorf("AdjustedRSquared() = %v should be less than R2 = %v", adjR2, r2)
	}
}

func TestErrorMetrics(t *testing.T) {
	actual := []float64{10, 20}
	predicted := []float64{11, 18}
	mae := MeanAbsoluteError(actual, predicted)
	if math.Abs(mae-1.5) > 0.0001 {
		t.Errorf("MeanAbsoluteError() = %v, want 1.5", mae)
	}
	rmse := RMSE(actual, predicted)
	if math.Abs(rmse-1.5811) > 0.001 {
		t.Errorf("RMSE() = %v, want ~1.5811", rmse)
	}
}
