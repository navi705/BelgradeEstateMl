package main

import (
	"sort"
)

func Avg(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range x {
		sum += value
	}
	return float64(sum) / float64(len(x))
}

func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func Pow(x float64, y int) float64 {
	if y == 0 {
		return 1
	}
	return x * Pow(x, y-1)
}

func Sqrt(x float64) float64 {
	if x < 0 || x == 0 {
		return 0
	}
	guess := x
	for {
		newGuess := (guess + x/guess) / 2
		if Abs(newGuess-guess) < 0.000001 {
			break
		}
		guess = newGuess
	}
	return guess
}

func Round(x float64, precision int) float64 {
	ratio := Pow(10, precision)
	return float64(int(x*ratio+0.5)) / ratio
}

func Correlation(x []float64, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}
	avgX := Avg(x)
	avgY := Avg(y)

	var numerator, sumSqX, sumSqY float64
	for i := 0; i < len(x); i++ {
		diffX := x[i] - avgX
		diffY := y[i] - avgY
		numerator += diffX * diffY
		sumSqX += diffX * diffX
		sumSqY += diffY * diffY
	}

	denominator := Sqrt(sumSqX * sumSqY)
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

func Median(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	sorted := make([]float64, len(x))
	copy(sorted, x)
	sort.Float64s(sorted)

	n := len(sorted)
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

func Mode(x []float64) []float64 {
	if len(x) == 0 {
		return nil
	}
	counts := make(map[float64]int)
	maxCount := 0
	for _, v := range x {
		counts[v]++
		if counts[v] > maxCount {
			maxCount = counts[v]
		}
	}

	var modes []float64
	for v, count := range counts {
		if count == maxCount {
			modes = append(modes, v)
		}
	}
	sort.Float64s(modes)
	return modes
}

func Percentile(x []float64, p float64) float64 {
	if len(x) == 0 {
		return 0
	}
	sorted := make([]float64, len(x))
	copy(sorted, x)
	sort.Float64s(sorted)

	if p <= 0 {
		return sorted[0]
	}
	if p >= 100 {
		return sorted[len(sorted)-1]
	}

	n := float64(len(sorted))
	index := (p / 100) * (n - 1)
	i := int(index)
	fraction := index - float64(i)

	if i+1 < len(sorted) {
		return sorted[i] + fraction*(sorted[i+1]-sorted[i])
	}
	return sorted[i]
}

func Quartile(x []float64, q int) float64 {
	switch q {
	case 1:
		return Percentile(x, 25)
	case 2:
		return Percentile(x, 50)
	case 3:
		return Percentile(x, 75)
	default:
		return 0
	}
}

func IQR(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	return Quartile(x, 3) - Quartile(x, 1)
}

func GetOutlierBounds(x []float64) (min float64, max float64) {
	if len(x) == 0 {
		return 0, 0
	}
	q1 := Quartile(x, 1)
	q3 := Quartile(x, 3)
	iqr := q3 - q1

	lower := q1 - 1.5*iqr
	upper := q3 + 1.5*iqr

	return lower, upper
}

func Histogram(data []float64, binCount int, rounded bool) []map[string]interface{} {
	if len(data) == 0 || binCount <= 0 {
		return nil
	}

	minVal := data[0]
	maxVal := data[0]
	for _, v := range data {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	if minVal == maxVal {
		resVal := minVal
		if rounded {
			resVal = float64(int(minVal + 0.5))
		}
		return []map[string]interface{}{
			{"from": resVal, "to": resVal, "count": len(data)},
		}
	}

	binWidth := (maxVal - minVal) / float64(binCount)
	bins := make([]map[string]interface{}, binCount)
	for i := 0; i < binCount; i++ {
		from := minVal + float64(i)*binWidth
		to := from + binWidth
		if i == binCount-1 {
			to = maxVal
		}

		fDisp := Round(from, 2)
		tDisp := Round(to, 2)
		if rounded {
			fDisp = float64(int(from + 0.5))
			tDisp = float64(int(to + 0.5))
		}

		bins[i] = map[string]interface{}{
			"from":  fDisp,
			"to":    tDisp,
			"count": 0,
		}
	}

	for _, v := range data {
		binIdx := int((v - minVal) / binWidth)
		if binIdx >= binCount {
			binIdx = binCount - 1
		}
		bins[binIdx]["count"] = bins[binIdx]["count"].(int) + 1
	}

	return bins
}

func Variance(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	avg := Avg(x)
	var sum float64
	for _, v := range x {
		diff := v - avg
		sum += diff * diff
	}
	return sum / float64(len(x))
}

func StdDev(x []float64) float64 {
	return Sqrt(Variance(x))
}

func IsNormalDistribution(x []float64) bool {
	if len(x) < 3 {
		return false
	}
	avg := Avg(x)
	sigma := StdDev(x)
	if sigma == 0 {
		return true
	}

	within1 := 0.0
	within2 := 0.0
	within3 := 0.0

	for _, v := range x {
		diff := Abs(v - avg)
		if diff <= sigma {
			within1++
		}
		if diff <= 2*sigma {
			within2++
		}
		if diff <= 3*sigma {
			within3++
		}
	}

	n := float64(len(x))
	p1 := within1 / n
	p2 := within2 / n
	p3 := within3 / n

	return p1 >= 0.6 && p1 <= 0.8 &&
		p2 >= 0.9 && p2 <= 1.0 &&
		p3 >= 0.98
}

func ZScore(x []float64) []float64 {
	if len(x) == 0 {
		return nil
	}
	mean := Avg(x)
	stdDev := StdDev(x)
	if stdDev == 0 {
		return make([]float64, len(x))
	}
	res := make([]float64, len(x))
	for i, v := range x {
		res[i] = (v - mean) / stdDev
	}
	return res
}

func LinearRegression(x, y []float64) (slope, intercept float64) {
	if len(x) != len(y) || len(x) < 2 {
		return 0, 0
	}
	avgX := Avg(x)
	avgY := Avg(y)
	var num, den float64
	for i := 0; i < len(x); i++ {
		num += (x[i] - avgX) * (y[i] - avgY)
		den += (x[i] - avgX) * (x[i] - avgX)
	}
	if den == 0 {
		return 0, 0
	}
	slope = num / den
	intercept = avgY - slope*avgX
	return slope, intercept
}

func GetSigmaBounds(x []float64, n float64) (min, max float64) {
	if len(x) == 0 {
		return 0, 0
	}
	avg := Avg(x)
	sigma := StdDev(x)
	return avg - n*sigma, avg + n*sigma
}

func RSquared(actual, predicted []float64) float64 {
	if len(actual) != len(predicted) || len(actual) == 0 {
		return 0
	}
	mean := Avg(actual)
	var ssRes, ssTot float64
	for i := 0; i < len(actual); i++ {
		ssRes += (actual[i] - predicted[i]) * (actual[i] - predicted[i])
		ssTot += (actual[i] - mean) * (actual[i] - mean)
	}
	if ssTot == 0 {
		return 1
	}
	return 1 - (ssRes / ssTot)
}

func SolveOLS(X [][]float64, Y []float64) []float64 {
	n := len(X)
	if n == 0 {
		return nil
	}
	m := len(X[0])
	XTX := make([][]float64, m)
	for i := 0; i < m; i++ {
		XTX[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			for k := 0; k < n; k++ {
				XTX[i][j] += X[k][i] * X[k][j]
			}
		}
	}
	XTY := make([]float64, m)
	for i := 0; i < m; i++ {
		for k := 0; k < n; k++ {
			XTY[i] += X[k][i] * Y[k]
		}
	}
	aug := make([][]float64, m)
	for i := 0; i < m; i++ {
		aug[i] = make([]float64, m+1)
		copy(aug[i], XTX[i])
		aug[i][m] = XTY[i]
	}
	for i := 0; i < m; i++ {
		pivot := i
		for j := i + 1; j < m; j++ {
			if Abs(aug[j][i]) > Abs(aug[pivot][i]) {
				pivot = j
			}
		}
		aug[i], aug[pivot] = aug[pivot], aug[i]
		if Abs(aug[i][i]) < 1e-10 {
			continue
		}
		for j := i + 1; j < m; j++ {
			factor := aug[j][i] / aug[i][i]
			for k := i; k <= m; k++ {
				aug[j][k] -= factor * aug[i][k]
			}
		}
	}
	weights := make([]float64, m)
	for i := m - 1; i >= 0; i-- {
		if Abs(aug[i][i]) < 1e-10 {
			weights[i] = 0
			continue
		}
		sum := aug[i][m]
		for j := i + 1; j < m; j++ {
			sum -= aug[i][j] * weights[j]
		}
		weights[i] = sum / aug[i][i]
	}
	return weights
}

func AdjustedRSquared(r2 float64, n, p int) float64 {
	if n <= p+1 {
		return 0
	}
	return 1 - (1-r2)*float64(n-1)/float64(n-p-1)
}

func MeanAbsoluteError(actual, predicted []float64) float64 {
	if len(actual) == 0 || len(actual) != len(predicted) {
		return 0
	}
	var sum float64
	for i := 0; i < len(actual); i++ {
		sum += Abs(actual[i] - predicted[i])
	}
	return sum / float64(len(actual))
}

func RMSE(actual, predicted []float64) float64 {
	if len(actual) == 0 || len(actual) != len(predicted) {
		return 0
	}
	var sum float64
	for i := 0; i < len(actual); i++ {
		diff := actual[i] - predicted[i]
		sum += diff * diff
	}
	return Sqrt(sum / float64(len(actual)))
}

func CrossValidate(X [][]float64, Y []float64, k int) float64 {
	n := len(X)
	if n < k || k <= 1 {
		return 0
	}

	foldSize := n / k
	totalR2 := 0.0

	for i := 0; i < k; i++ {
		var trainX [][]float64
		var trainY []float64
		var testX [][]float64
		var testY []float64

		start := i * foldSize
		end := start + foldSize
		if i == k-1 {
			end = n
		}

		for j := 0; j < n; j++ {
			if j >= start && j < end {
				testX = append(testX, X[j])
				testY = append(testY, Y[j])
			} else {
				trainX = append(trainX, X[j])
				trainY = append(trainY, Y[j])
			}
		}

		weights := SolveOLS(trainX, trainY)
		if weights == nil {
			continue
		}

		predicted := make([]float64, len(testY))
		for j := range testY {
			val := 0.0
			for l := range weights {
				val += testX[j][l] * weights[l]
			}
			predicted[j] = val
		}

		totalR2 += RSquared(testY, predicted)
	}

	return totalR2 / float64(k)
}

func EuclideanDistance(v1, v2 []float64) float64 {
	if len(v1) != len(v2) || len(v1) == 0 {
		return 0
	}
	var sum float64
	for i := range v1 {
		diff := v1[i] - v2[i]
		sum += diff * diff
	}
	return Sqrt(sum)
}
