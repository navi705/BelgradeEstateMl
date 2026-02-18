package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "realestate_api_requests_total",
		Help: "Total number of API requests",
	}, []string{"endpoint", "district"})

	metricPredictionPrice = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "realestate_prediction_price_last",
		Help: "Last predicted price for a district and algorithm",
	}, []string{"algorithm", "district"})

	metricModelR2 = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "realestate_model_r2",
		Help: "Model R-squared for fit quality",
	}, []string{"algorithm", "district"})

	metricModelMAE = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "realestate_model_mae",
		Help: "Model Mean Absolute Error",
	}, []string{"algorithm", "district"})

	metricMarketTrend = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "realestate_market_trend_percent",
		Help: "Calculated market trend percentage",
	}, []string{"district"})

	metricPredictionSqm = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "realestate_prediction_sqm_last",
		Help: "SQM of the last prediction",
	}, []string{"algorithm", "district"})

	metricPredictionRooms = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "realestate_prediction_rooms_last",
		Help: "Rooms of the last prediction",
	}, []string{"algorithm", "district"})

	metricPredictionFloor = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "realestate_prediction_floor_last",
		Help: "Floor of the last prediction",
	}, []string{"algorithm", "district"})
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	storage, err := NewConnection(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer storage.db.Close()

	// mux := http.NewServeMux() // No longer needed if using default ServeMux

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		discovery := map[string]interface{}{
			"name":        "BelgradeEstateML API",
			"description": "Professional real estate analytics and predictive engine for Belgrade",
			"version":     "2.0",
			"endpoints": []map[string]interface{}{
				{"path": "/", "description": "API Discovery (this page)"},
				{"path": "/districts", "description": "List all available municipalities for filtering"},
				{"path": "/correlation", "description": "Feature correlation matrix", "params": []string{"from", "to", "district", "round"}},
				{"path": "/stats", "description": "Basic statistics for a field", "params": []string{"field", "from", "to", "district", "round"}},
				{"path": "/analyze", "description": "Advanced analytics with normality and outlier detection", "params": []string{"fields", "outlier_method", "outlier_field", "from", "to", "district", "round"}},
				{"path": "/predict", "description": "Linear/Polynomial price prediction with diagnostics", "params": []string{"sqm", "rooms", "floor", "district", "round"}},
				{"path": "/predict/knn", "description": "K-Nearest Neighbors price prediction", "params": []string{"sqm", "rooms", "floor", "district", "round"}},
				{"path": "/predict/tree", "description": "Decision Tree price prediction", "params": []string{"sqm", "rooms", "floor", "district", "round"}},
				{"path": "/predict/boost", "description": "Gradient Boosting (Ensemble) price prediction", "params": []string{"sqm", "rooms", "floor", "district", "round"}},
			},
			"example": "/predict?district=Vracar&sqm=60&rooms=2&floor=3",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(discovery)
	})

	getRoundParam := func(r *http.Request, defaultVal int) int {
		rStr := r.URL.Query().Get("round")
		if rStr == "" {
			return defaultVal
		}
		val, err := strconv.Atoi(rStr)
		if err != nil {
			return defaultVal
		}
		return val
	}

	getFilteredData := func(r *http.Request) ([]RealEstate, time.Time, time.Time, string, error) {
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		district := r.URL.Query().Get("district")
		if district != "" {
			district = StandardizeDistrict(district)
		}
		excludeOutliers := r.URL.Query().Get("exclude_outliers") != "false"

		var from, to time.Time
		if fromStr != "" {
			from, _ = time.Parse("2006-01-02", fromStr)
		}
		if toStr != "" {
			to, _ = time.Parse("2006-01-02", toStr)
		}

		if from.IsZero() || to.IsZero() {
			min, max, _ := GetDateRange(storage)
			if from.IsZero() {
				from = min
			}
			if to.IsZero() {
				to = max
			}
		}

		estates, err := GetRealEstateWithoutDuplicate(storage, from, to)
		if err != nil {
			return nil, from, to, district, err
		}

		if district != "" {
			estates = FilterByDistrict(estates, district)
		}

		if excludeOutliers {
			method := r.URL.Query().Get("outlier_method")
			estates = AggressiveClean(estates, method)
		}

		return estates, from, to, district, nil
	}

	calculateFieldStats := func(data []float64, rounded bool, precision int) map[string]interface{} {
		if len(data) == 0 {
			return nil
		}
		return map[string]interface{}{
			"avg":          Round(Avg(data), precision),
			"median":       Round(Median(data), precision),
			"mode":         Mode(data),
			"q1":           Round(Quartile(data, 1), precision),
			"q3":           Round(Quartile(data, 3), precision),
			"is_normal":    IsNormalDistribution(data),
			"distribution": Histogram(data, 10, rounded),
		}
	}

	generateAllStats := func(estates []RealEstate, precision int) map[string]interface{} {
		if len(estates) == 0 {
			return nil
		}
		fields := map[string][]float64{
			"price":       make([]float64, len(estates)),
			"sqm":         make([]float64, len(estates)),
			"rooms":       make([]float64, len(estates)),
			"floor":       make([]float64, len(estates)),
			"floor_total": make([]float64, len(estates)),
		}

		for i, e := range estates {
			fields["price"][i] = float64(e.Price)
			fields["sqm"][i] = float64(e.SquareMeter)
			fields["rooms"][i] = float64(e.QuantityRoom)
			fields["floor"][i] = float64(e.Floor)
			fields["floor_total"][i] = float64(e.FloorTotal)
		}

		res := make(map[string]interface{})
		for name, data := range fields {
			rounded := (name == "rooms" || name == "floor" || name == "floor_total")
			effPrecision := precision
			if rounded {
				effPrecision = 1
			}
			res[name] = calculateFieldStats(data, rounded, effPrecision)
		}
		return res
	}

	http.HandleFunc("/full", func(w http.ResponseWriter, r *http.Request) {
		estates, from, to, district, err := getFilteredData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		precision := getRoundParam(r, 2)
		matrix := CorrelationMatrix(estates)
		stats := generateAllStats(estates, precision)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"district":    district,
			"from":        from.Format("2006-01-02"),
			"to":          to.Format("2006-01-02"),
			"count":       len(estates),
			"stats":       stats,
			"correlation": matrix,
			"data":        estates,
		})
	})

	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		district := r.URL.Query().Get("district")
		if district != "" {
			district = StandardizeDistrict(district)
		}

		method := r.URL.Query().Get("outlier_method")
		field := r.URL.Query().Get("outlier_field")
		fieldsStr := r.URL.Query().Get("fields")

		var from, to time.Time
		if fromStr != "" {
			from, _ = time.Parse("2006-01-02", fromStr)
		}
		if toStr != "" {
			to, _ = time.Parse("2006-01-02", toStr)
		}

		if from.IsZero() || to.IsZero() {
			min, max, _ := GetDateRange(storage)
			if from.IsZero() {
				from = min
			}
			if to.IsZero() {
				to = max
			}
		}

		estates, err := GetRealEstateWithoutDuplicate(storage, from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if district != "" {
			estates = FilterByDistrict(estates, district)
		}

		if method != "" {
			estates = FilterOutliersConfigurable(estates, field, method)
		}

		precision := getRoundParam(r, 2)
		allStats := generateAllStats(estates, precision)
		var resultStats map[string]interface{}

		if fieldsStr != "" {
			resultStats = make(map[string]interface{})
			requested := strings.Split(fieldsStr, ",")
			for _, f := range requested {
				f = strings.TrimSpace(f)
				if s, ok := allStats[f]; ok {
					resultStats[f] = s
				}
			}
		} else {
			resultStats = allStats
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"district":       district,
			"from":           from.Format("2006-01-02"),
			"to":             to.Format("2006-01-02"),
			"count":          len(estates),
			"outlier_method": method,
			"outlier_field":  field,
			"monthly_trend":  Round(CalculateTrend(estates), precision),
			"stats":          resultStats,
		})
	})

	http.HandleFunc("/correlation", func(w http.ResponseWriter, r *http.Request) {
		estates, from, to, district, err := getFilteredData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		precision := getRoundParam(r, 2)
		matrix := CorrelationMatrix(estates)
		for i := range matrix {
			for j := range matrix[i] {
				matrix[i][j] = Round(matrix[i][j], precision)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"district":    district,
			"from":        from.Format("2006-01-02"),
			"to":          to.Format("2006-01-02"),
			"count":       len(estates),
			"labels":      []string{"Price", "Sqm", "Rooms", "Floor", "FloorTotal"},
			"correlation": matrix,
		})
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		estates, from, to, district, err := getFilteredData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(estates) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "no data for current filters"})
			return
		}

		precision := getRoundParam(r, 2)
		stats := generateAllStats(estates, precision)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"district": district,
			"from":     from.Format("2006-01-02"),
			"to":       to.Format("2006-01-02"),
			"count":    len(estates),
			"stats":    stats,
		})
	})

	http.HandleFunc("/districts", func(w http.ResponseWriter, r *http.Request) {
		districts := GetAllStandardizedDistricts()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(districts)
	})

	http.HandleFunc("/period", func(w http.ResponseWriter, r *http.Request) {
		min, max, err := GetDateRange(storage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"from": min.Format("2006-01-02"),
			"to":   max.Format("2006-01-02"),
		})
	})

	http.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		district := r.URL.Query().Get("district")
		if district != "" {
			district = StandardizeDistrict(district)
		}

		sqm, _ := strconv.ParseFloat(r.URL.Query().Get("sqm"), 64)
		rooms, _ := strconv.ParseFloat(r.URL.Query().Get("rooms"), 64)
		floor, _ := strconv.ParseFloat(r.URL.Query().Get("floor"), 64)

		estates, err := GetRealEstateWithoutDuplicate(storage, time.Time{}, time.Time{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if district != "" {
			estates = FilterByDistrict(estates, district)
		}

		method := r.URL.Query().Get("outlier_method")
		estates = AggressiveClean(estates, method)

		precision := getRoundParam(r, 0)
		model := TrainModel(estates)
		pred, pMin, pMax := model.PredictWithInterval(sqm, rooms, floor)

		metricRequests.WithLabelValues("/predict", district).Inc()
		metricPredictionPrice.WithLabelValues("polynomial", district).Set(pred)
		metricPredictionSqm.WithLabelValues("polynomial", district).Set(sqm)
		metricPredictionRooms.WithLabelValues("polynomial", district).Set(rooms)
		metricPredictionFloor.WithLabelValues("polynomial", district).Set(floor)
		metricModelR2.WithLabelValues("polynomial", district).Set(model.RSquared)
		metricModelMAE.WithLabelValues("polynomial", district).Set(model.MAE)
		metricMarketTrend.WithLabelValues(district).Set(model.Trend)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"district":    district,
			"sqm":         sqm,
			"rooms":       rooms,
			"floor":       floor,
			"prediction":  Round(pred, precision),
			"price_min":   Round(pMin, precision),
			"price_max":   Round(pMax, precision),
			"r2":          Round(model.RSquared, 4),
			"adjusted_r2": Round(model.AdjustedR2, 4),
			"cv_score":    Round(model.CVScore, 4),
			"mae":         Round(model.MAE, precision),
			"rmse":        Round(model.RMSE, precision),
			"trend":       Round(model.Trend, 2),
			"status":      model.Status,
			"condition":   model.Condition,
			"count":       model.Count,
		})
	})

	http.HandleFunc("/predict/knn", func(w http.ResponseWriter, r *http.Request) {
		estates, _, _, district, err := getFilteredData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sqm, _ := strconv.ParseFloat(r.URL.Query().Get("sqm"), 64)
		rooms, _ := strconv.ParseFloat(r.URL.Query().Get("rooms"), 64)
		floor, _ := strconv.ParseFloat(r.URL.Query().Get("floor"), 64)

		precision := getRoundParam(r, 0)
		prediction := PredictKNN(estates, sqm, rooms, floor, 10)

		metricRequests.WithLabelValues("/predict/knn", district).Inc()
		metricPredictionPrice.WithLabelValues("knn", district).Set(prediction)
		metricPredictionSqm.WithLabelValues("knn", district).Set(sqm)
		metricPredictionRooms.WithLabelValues("knn", district).Set(rooms)
		metricPredictionFloor.WithLabelValues("knn", district).Set(floor)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"prediction": Round(prediction, precision),
			"algorithm":  "KNN",
			"k":          10,
			"count":      len(estates),
		})
	})

	http.HandleFunc("/predict/tree", func(w http.ResponseWriter, r *http.Request) {
		estates, _, _, district, err := getFilteredData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		X := make([][]float64, len(estates))
		Y := make([]float64, len(estates))
		for i, e := range estates {
			X[i] = []float64{float64(e.SquareMeter), float64(e.QuantityRoom), float64(e.Floor)}
			Y[i] = float64(e.Price)
		}

		tree := BuildTree(X, Y, 0, 5)
		sqm, _ := strconv.ParseFloat(r.URL.Query().Get("sqm"), 64)
		rooms, _ := strconv.ParseFloat(r.URL.Query().Get("rooms"), 64)
		floor, _ := strconv.ParseFloat(r.URL.Query().Get("floor"), 64)

		precision := getRoundParam(r, 0)
		prediction := 0.0
		if tree != nil {
			prediction = tree.Predict([]float64{sqm, rooms, floor})
		}

		metricRequests.WithLabelValues("/predict/tree", district).Inc()
		metricPredictionPrice.WithLabelValues("tree", district).Set(prediction)
		metricPredictionSqm.WithLabelValues("tree", district).Set(sqm)
		metricPredictionRooms.WithLabelValues("tree", district).Set(rooms)
		metricPredictionFloor.WithLabelValues("tree", district).Set(floor)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"prediction": Round(prediction, precision),
			"algorithm":  "Decision Tree",
			"max_depth":  5,
			"count":      len(estates),
		})
	})

	http.HandleFunc("/predict/boost", func(w http.ResponseWriter, r *http.Request) {
		estates, _, _, district, err := getFilteredData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		X := make([][]float64, len(estates))
		Y := make([]float64, len(estates))
		for i, e := range estates {
			X[i] = []float64{float64(e.SquareMeter), float64(e.QuantityRoom), float64(e.Floor)}
			Y[i] = float64(e.Price)
		}

		model := TrainBoosting(X, Y, 20, 0.1)
		sqm, _ := strconv.ParseFloat(r.URL.Query().Get("sqm"), 64)
		rooms, _ := strconv.ParseFloat(r.URL.Query().Get("rooms"), 64)
		floor, _ := strconv.ParseFloat(r.URL.Query().Get("floor"), 64)

		precision := getRoundParam(r, 0)
		prediction := 0.0
		if model != nil {
			prediction = model.Predict([]float64{sqm, rooms, floor})
		}

		metricRequests.WithLabelValues("/predict/boost", district).Inc()
		metricPredictionPrice.WithLabelValues("boost", district).Set(prediction)
		metricPredictionSqm.WithLabelValues("boost", district).Set(sqm)
		metricPredictionRooms.WithLabelValues("boost", district).Set(rooms)
		metricPredictionFloor.WithLabelValues("boost", district).Set(floor)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"prediction":    Round(prediction, precision),
			"algorithm":     "Gradient Boosting",
			"trees":         20,
			"learning_rate": 0.1,
			"count":         len(estates),
		})
	})

	http.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting ML server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
