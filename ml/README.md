## 🧠 Prediction Methodology

The system uses a multi-model approach to valuation. Each endpoint represents a different mathematical perspective:

### 1. Standard Regression (`/predict`)
- **Type**: Linear OLS (Ordinary Least Squares) with Polynomial Features.
- **How it works**: Models the price as a continuous function of area, rooms, and floor. It also captures non-linear price growth (e.g., how sqm affects price exponentially in luxury areas).
- **Diagnostics**: Returns R² (fit quality), MAE (average error), and a **Price corridor** (min/max price bounds).

### 2. K-Nearest Neighbors (`/predict/knn`)
- **Type**: Distance-based non-parametric model.
- **How it works**: Finds the 10 apartments most similar to your request in terms of location and size, then averages their prices.
- **Best Use**: Most "human-like" valuation; mimics the logic of real estate agents comparing similar objects.

### 3. Decision Tree (`/predict/tree`)
- **Type**: Recursive Binary Partitioning.
- **How it works**: Splits the entire market data into branches (e.g., "Is it in Vracar?" -> "Is it above 5th floor?").
- **Best Use**: Capturing sharp price differences based on specific thresholds (e.g., radical price jump for ground floor vs first floor).

### 4. Gradient Boosting (`/predict/boost`)
- **Type**: Ensemble Learning (Gradient Boosted Regression Trees).
- **How it works**: Starts with a simple tree and builds 20 subsequent trees, where each tree specifically tries to correct the errors made by all previous ones.
- **Best Use**: Maximum accuracy. This is the "gold standard" for real estate tabular data.

---

## 📊 Analytics & Market Insights

Beyond simple valuation, you can analyze market health:

- **`/analyze`**: Performs normality tests to see if market prices are "natural" or manipulated. Automatically identifies and removes "statistical noise" (outliers).
- **`/correlation`**: See how strongly parameters are linked. (Does floor really affect price in Novi Beograd? Check the coefficient here).
- **Trends**: All prediction outputs include a `monthly_trend` showing the % change in price over time for that specific area.

---

## 🛠 Query Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `district` | String | Municipality name (e.g., `Zemun`, `Vracar`). Case-insensitive. |
| `sqm` | Float | Living area in square meters. |
| `rooms` | Float | Number of rooms. |
| `floor` | Float | Floor number. |
| `from` / `to` | Date | Filter data by date (`YYYY-MM-DD`). |
| `round` | Int | Control response precision (e.g., `round=0` for whole integers). |
| `outlier_method` | String | `sigma` (3-sigma rule) or `iqr` (interquartile range). |
| `exclude_outliers` | Bool | Set to `false` to include outliers (Defaults to `true` for all analytics and predictions). |

---

---

## 📖 API Reference & Examples

### 1. API Discovery
**Endpoint:** `GET /`
**Description:** Lists all available endpoints.
**Request:** `GET /`
**Response:**
```json
{
  "name": "BelgradeEstateML API",
  "version": "2.0",
  "endpoints": [
    {"path": "/", "description": "API Discovery (this page)"},
    {"path": "/predict", "description": "Linear/Polynomial price prediction", "params": ["sqm", "rooms", "floor", "district", "round"]}
  ],
  "example": "/predict?district=Vracar&sqm=60&rooms=2&floor=3"
}
```

### 2. Standard Valuation (Linear)
**Endpoint:** `GET /predict`
**Request:** `GET /predict?district=Zemun&sqm=50&rooms=2&floor=4&round=0`
**Response:**
```json
{
  "prediction": 112000,
  "price_min": 108500,
  "price_max": 115500,
  "status": 1,
  "condition": "Success",
  "r2": 0.8542,
  "mae": 3500,
  "trend": 1.2
}
```

### 3. AI Valuation (Gradient Boosting)
**Endpoint:** `GET /predict/boost`
**Request:** `GET /predict/boost?district=Vracar&sqm=80&rooms=3&floor=2&round=0`
**Response:**
```json
{
  "prediction": 245000,
  "algorithm": "Gradient Boosting",
  "trees": 20,
  "count": 450
}
```

### 4. Advanced Analytics
**Endpoint:** `GET /analyze`
**Request:** `GET /analyze?district=Novi+Beograd&fields=price&outlier_method=iqr&round=2`
**Response:**
```json
{
  "district": "Novi Beograd",
  "count": 850,
  "monthly_trend": 0.85,
  "stats": {
    "price": {
      "avg": 2150.45,
      "median": 2100.0,
      "is_normal": true,
      "distribution": [
        {"from": 1500, "to": 1700, "count": 45},
        {"from": 1700, "to": 1900, "count": 120}
      ]
    }
  }
}
```

### 5. Feature Correlation
**Endpoint:** `GET /correlation`
**Request:** `GET /correlation?district=Stari+Grad&round=3`
**Response:**
```json
{
  "district": "Stari Grad",
  "labels": ["Price", "Sqm", "Rooms", "Floor", "FloorTotal"],
  "correlation": [
    [1.0, 0.92, 0.75, 0.12, 0.05],
    [0.92, 1.0, 0.81, 0.15, 0.08]
  ]
}
```

### 6. Available Districts
**Endpoint:** `GET /districts`
**Request:** `GET /districts`
**Response:**
```json
["Vracar", "Novi Beograd", "Zemun", "Palilula", "Zvezdara", "Stari Grad", "Savski Venac", "Voždovac"]
```

### 7. Data Availability Period
**Endpoint:** `GET /period`
**Request:** `GET /period`
**Response:**
```json
{
  "from": "2023-10-01",
  "to": "2024-02-18"
}
```
