package main

import (
	"sort"
	"strings"
)

func AggressiveClean(estates []RealEstate, method string) []RealEstate {
	if len(estates) < 10 {
		return estates
	}

	if method == "" {
		method = "iqr"
	}

	filtered := FilterOutliersConfigurable(estates, "price", method)
	if len(filtered) > 10 {
		filtered = FilterOutliersConfigurable(filtered, "sqm", method)
	}

	if len(filtered) > 10 {
		var ppsqms []float64
		for _, e := range filtered {
			if e.SquareMeter > 0 {
				ppsqms = append(ppsqms, float64(e.Price)/float64(e.SquareMeter))
			}
		}

		var lower, upper float64
		if method == "sigma" {
			lower, upper = GetSigmaBounds(ppsqms, 3)
		} else {
			lower, upper = GetOutlierBounds(ppsqms)
		}

		var result []RealEstate
		for _, e := range filtered {
			val := float64(e.Price) / float64(e.SquareMeter)
			if val >= lower && val <= upper {
				result = append(result, e)
			}
		}
		filtered = result
	}

	return filtered
}

func CorrelationMatrix(estates []RealEstate) [][]float64 {
	if len(estates) == 0 {
		return nil
	}

	props := [][]float64{
		make([]float64, len(estates)),
		make([]float64, len(estates)),
		make([]float64, len(estates)),
		make([]float64, len(estates)),
		make([]float64, len(estates)),
	}

	for i, e := range estates {
		props[0][i] = float64(e.Price)
		props[1][i] = float64(e.SquareMeter)
		props[2][i] = float64(e.QuantityRoom)
		props[3][i] = float64(e.Floor)
		props[4][i] = float64(e.FloorTotal)
	}

	matrix := make([][]float64, len(props))
	for i := 0; i < len(props); i++ {
		matrix[i] = make([]float64, len(props))
		for j := 0; j < len(props); j++ {
			matrix[i][j] = Correlation(props[i], props[j])
		}
	}
	return matrix
}

func FilterByDistrict(estates []RealEstate, district string) []RealEstate {
	var filtered []RealEstate
	for _, e := range estates {
		if strings.EqualFold(e.District, district) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func FilterOutliers(estates []RealEstate) []RealEstate {
	if len(estates) < 4 {
		return estates
	}

	prices := make([]float64, len(estates))
	for i, e := range estates {
		prices[i] = float64(e.Price)
	}

	lower, upper := GetOutlierBounds(prices)

	var filtered []RealEstate
	for _, e := range estates {
		p := float64(e.Price)
		if p >= lower && p <= upper {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

func FilterOutliersConfigurable(estates []RealEstate, field string, method string) []RealEstate {
	if len(estates) < 4 {
		return estates
	}

	values := make([]float64, len(estates))
	for i, e := range estates {
		switch field {
		case "sqm":
			values[i] = float64(e.SquareMeter)
		case "rooms":
			values[i] = float64(e.QuantityRoom)
		case "floor":
			values[i] = float64(e.Floor)
		case "floor_total":
			values[i] = float64(e.FloorTotal)
		default:
			values[i] = float64(e.Price)
		}
	}

	var lower, upper float64
	if method == "sigma" {
		lower, upper = GetSigmaBounds(values, 3)
	} else {
		lower, upper = GetOutlierBounds(values)
	}

	var filtered []RealEstate
	for _, e := range estates {
		var val float64
		switch field {
		case "sqm":
			val = float64(e.SquareMeter)
		case "rooms":
			val = float64(e.QuantityRoom)
		case "floor":
			val = float64(e.Floor)
		case "floor_total":
			val = float64(e.FloorTotal)
		default:
			val = float64(e.Price)
		}

		if val >= lower && val <= upper {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

type PredictiveModel struct {
	Weights    []float64
	RSquared   float64
	AdjustedR2 float64
	MAE        float64
	RMSE       float64
	CVScore    float64
	Trend      float64
	Status     int
	Condition  string
	Count      int
}

func TrainModel(estates []RealEstate) PredictiveModel {
	if len(estates) < 4 {
		return PredictiveModel{Status: 4, Condition: "Insufficient Data", Count: len(estates)}
	}

	X := make([][]float64, len(estates))
	Y := make([]float64, len(estates))
	for i, e := range estates {
		sqm := float64(e.SquareMeter)
		X[i] = []float64{1, sqm, sqm * sqm, float64(e.QuantityRoom), float64(e.Floor)}
		Y[i] = float64(e.Price)
	}

	weights := SolveOLS(X, Y)
	if weights == nil {
		return PredictiveModel{Status: 3, Condition: "Meaningless/Error", Count: len(estates)}
	}

	cvScore := CrossValidate(X, Y, 5)

	predicted := make([]float64, len(estates))
	for i := range estates {
		val := 0.0
		for j := range weights {
			val += X[i][j] * weights[j]
		}
		predicted[i] = val
	}

	r2 := RSquared(Y, predicted)
	count := len(estates)
	adjR2 := AdjustedRSquared(r2, count, 4)
	mae := MeanAbsoluteError(Y, predicted)
	rmse := RMSE(Y, predicted)
	trend := CalculateTrend(estates)

	status := 3
	condition := "Unreliable/Poor Fit"

	if r2 > 0.6 && count >= 20 {
		status = 1
		condition = "Success"
	}
	if r2 > 0.9 && count < 15 {
		status = 2
		condition = "Potential Overfit"
	}
	if count < 10 {
		status = 4
		condition = "Insufficient Data"
	}

	return PredictiveModel{
		Weights:    weights,
		RSquared:   r2,
		AdjustedR2: adjR2,
		MAE:        mae,
		RMSE:       rmse,
		CVScore:    cvScore,
		Trend:      trend,
		Status:     status,
		Condition:  condition,
		Count:      count,
	}
}

func (m PredictiveModel) Predict(sqm, rooms, floor float64) float64 {
	if len(m.Weights) < 5 {
		return 0
	}
	return m.Weights[0] + sqm*m.Weights[1] + (sqm*sqm)*m.Weights[2] + rooms*m.Weights[3] + floor*m.Weights[4]
}

func CalculateTrend(estates []RealEstate) float64 {
	if len(estates) < 10 {
		return 0
	}

	monthlyPrices := make(map[string][]float64)
	for _, e := range estates {
		month := e.ParsingDate.Format("2006-01")
		if e.SquareMeter > 0 {
			sqmPrice := float64(e.Price) / float64(e.SquareMeter)
			monthlyPrices[month] = append(monthlyPrices[month], sqmPrice)
		}
	}

	var months []string
	for m := range monthlyPrices {
		months = append(months, m)
	}
	sort.Strings(months)

	if len(months) < 2 {
		return 0
	}

	avgPrices := make([]float64, len(months))
	for i, m := range months {
		avgPrices[i] = Avg(monthlyPrices[m])
	}

	var totalChange float64
	for i := 1; i < len(avgPrices); i++ {
		if avgPrices[i-1] != 0 {
			change := (avgPrices[i] - avgPrices[i-1]) / avgPrices[i-1]
			totalChange += change
		}
	}

	return (totalChange / float64(len(avgPrices)-1)) * 100
}

func (m PredictiveModel) PredictWithInterval(sqm, rooms, floor float64) (price, min, max float64) {
	price = m.Predict(sqm, rooms, floor)
	if price <= 0 {
		return 0, 0, 0
	}

	spread := m.MAE
	if m.Status == 2 || m.Status == 3 {
		spread *= 1.5
	}

	min = price - spread
	max = price + spread

	if min < 0 {
		min = 0
	}

	return price, min, max
}

func PredictKNN(estates []RealEstate, targetSqm, targetRooms, targetFloor float64, k int) float64 {
	if len(estates) == 0 || k <= 0 {
		return 0
	}

	type neighbor struct {
		distance float64
		price    float64
	}

	neighbors := make([]neighbor, len(estates))
	for i, e := range estates {
		// Features: sqm, rooms, floor
		v1 := []float64{targetSqm, targetRooms, targetFloor}
		v2 := []float64{float64(e.SquareMeter), float64(e.QuantityRoom), float64(e.Floor)}
		neighbors[i] = neighbor{
			distance: EuclideanDistance(v1, v2),
			price:    float64(e.Price),
		}
	}

	sort.Slice(neighbors, func(i, j int) bool {
		return neighbors[i].distance < neighbors[j].distance
	})

	if k > len(neighbors) {
		k = len(neighbors)
	}

	var sum float64
	for i := 0; i < k; i++ {
		sum += neighbors[i].price
	}
	return sum / float64(k)
}

type Node struct {
	FeatureIndex int
	Threshold    float64
	Value        float64
	Left         *Node
	Right        *Node
}

func BuildTree(X [][]float64, Y []float64, depth, maxDepth int) *Node {
	if len(Y) == 0 {
		return nil
	}

	if depth >= maxDepth || len(Y) < 5 {
		return &Node{Value: Avg(Y)}
	}

	bestFeature := -1
	bestThreshold := 0.0
	minVariance := -1.0

	for f := 0; f < len(X[0]); f++ {
		for _, row := range X {
			threshold := row[f]
			var leftY, rightY []float64
			for i, v := range Y {
				if X[i][f] <= threshold {
					leftY = append(leftY, v)
				} else {
					rightY = append(rightY, v)
				}
			}

			if len(leftY) == 0 || len(rightY) == 0 {
				continue
			}

			variance := (Variance(leftY) * float64(len(leftY))) + (Variance(rightY) * float64(len(rightY)))
			if minVariance == -1 || variance < minVariance {
				minVariance = variance
				bestFeature = f
				bestThreshold = threshold
			}
		}
	}

	if bestFeature == -1 {
		return &Node{Value: Avg(Y)}
	}

	var leftX, rightX [][]float64
	var leftY, rightY []float64
	for i, v := range Y {
		if X[i][bestFeature] <= bestThreshold {
			leftX = append(leftX, X[i])
			leftY = append(leftY, v)
		} else {
			rightX = append(rightX, X[i])
			rightY = append(rightY, v)
		}
	}

	return &Node{
		FeatureIndex: bestFeature,
		Threshold:    bestThreshold,
		Left:         BuildTree(leftX, leftY, depth+1, maxDepth),
		Right:        BuildTree(rightX, rightY, depth+1, maxDepth),
	}
}

func (n *Node) Predict(features []float64) float64 {
	if n.Left == nil && n.Right == nil {
		return n.Value
	}
	if features[n.FeatureIndex] <= n.Threshold {
		return n.Left.Predict(features)
	}
	return n.Right.Predict(features)
}

type BoostingModel struct {
	Trees        []*Node
	LearningRate float64
}

func TrainBoosting(X [][]float64, Y []float64, nTrees int, lr float64) *BoostingModel {
	if len(Y) == 0 {
		return nil
	}

	model := &BoostingModel{LearningRate: lr}
	residuals := make([]float64, len(Y))
	copy(residuals, Y)

	for i := 0; i < nTrees; i++ {
		tree := BuildTree(X, residuals, 0, 3)
		if tree == nil {
			break
		}
		model.Trees = append(model.Trees, tree)

		for j := range residuals {
			residuals[j] -= lr * tree.Predict(X[j])
		}
	}
	return model
}

func (m *BoostingModel) Predict(features []float64) float64 {
	var prediction float64
	for _, tree := range m.Trees {
		prediction += m.LearningRate * tree.Predict(features)
	}
	return prediction
}
