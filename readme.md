# Belgrade Estate ML

Analytical engine and predictive models for the Belgrade real estate market. The system processes parsed data to provide valuation estimates, market trends, and statistical insights.

## API Reference

Base URL: [https://bg-real-estate.duckdns.org](https://bg-real-estate.duckdns.org)

### Municipalities
Get a list of all supported districts for filtering.
- **URL**: [https://bg-real-estate.duckdns.org/districts](https://bg-real-estate.duckdns.org/districts)

### Valuation and Predictions
Predict apartment prices using different mathematical models.

**Parameters:**
- `sqm`: Square meters (required)
- `rooms`: Number of rooms
- `floor`: Floor number
- `district`: Municipality name (e.g., Vracar, Zemun)
- `round`: Decimal precision (0-4)
- `outlier_method`: Method for outlier removal (`iqr` or `sigma`). **IQR is recommended**.
- `exclude_outliers`: Enable/disable outlier filtering (`true` or `false`).

> [!TIP]
> **Use IQR**: The `sigma` method assumes a normal distribution. Real estate prices are often skewed. If `is_normal: false` in stats, always use **IQR**.

**Gradient Boosting (Recommended):**
- [Test Link: 65m² in Vracar (Sigma)](https://bg-real-estate.duckdns.org/predict/boost?sqm=65&rooms=2.5&floor=3&district=Vracar&outlier_method=sigma)
- [Test Link: 65m² in Vracar (IQR)](https://bg-real-estate.duckdns.org/predict/boost?sqm=65&rooms=2.5&floor=3&district=Vracar&outlier_method=iqr)

**K-Nearest Neighbors:**
- [Test Link: 40m² in Zemun](https://bg-real-estate.duckdns.org/predict/knn?sqm=40&rooms=1&district=Zemun&exclude_outliers=true&outlier_method=iqr)

**Decision Tree:**
- [Test Link: 120m² in Savski Venac](https://bg-real-estate.duckdns.org/predict/tree?sqm=120&rooms=4&district=Savski+Venac)

**Polynomial Regression:**
- [Test Link: 55m² in Zvezdara](https://bg-real-estate.duckdns.org/predict?sqm=55&district=Zvezdara&round=2)

### Analytics and Statistics

**Advanced Analysis & Market Trends:**
Returns normality checks, distribution data, and the **monthly market trend**.
- [Test Link: Trends & Stats for Palilula](https://bg-real-estate.duckdns.org/analyze?fields=price&district=Palilula&outlier_method=iqr)

**Feature Correlation:**
Matrix of how price, area, and rooms relate to each other.
- [Test Link: Correlation in Vozdovac (2024)](https://bg-real-estate.duckdns.org/correlation?district=Vozdovac&from=2024-01-01&to=2024-12-31)

**Field Statistics:**
Basic stats (min, max, mean, median) for a specific metric.
- [Test Link: Statistics for Cukarica](https://bg-real-estate.duckdns.org/stats?district=Cukarica&exclude_outliers=true)
