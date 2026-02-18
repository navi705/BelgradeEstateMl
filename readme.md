# Belgrade Estate ML

Analytical engine and predictive models for the Belgrade real estate market. The system processes parsed data to provide valuation estimates, market trends, and statistical insights.

## API Reference

Base URL: [https://bg-real-estate.duckdns.org](https://bg-real-estate.duckdns.org)

### System Info & Discovery
- **Districts**: Get all supported municipalities.
  - [Test Link](https://bg-real-estate.duckdns.org/districts)
- **Data Period**: Get the available date range (min/max date) in the database.
  - [Test Link](https://bg-real-estate.duckdns.org/period)
- **Rate Limit**: 2 requests/sec per IP (Burst: 5).

### Valuation and Predictions
Predict apartment prices using different mathematical models.

**Parameters:**
- `sqm`: Square meters (required)
- `rooms`: Number of rooms
- `floor`: Floor number
- `district`: Municipality name (e.g., Vračar, Zemun)
- `from`: Start date for data filtering (`YYYY-MM-DD`)
- `to`: End date for data filtering (`YYYY-MM-DD`)
- `round`: Decimal precision (0-4)
- `outlier_method`: Method for outlier removal (`iqr` or `sigma`). **IQR is recommended**.
- `exclude_outliers`: Enable/disable outlier filtering (`true` or `false`).

> [!TIP]
> **Use IQR**: The `sigma` method assumes a normal distribution. Real estate prices are often skewed. If `is_normal: false` in stats, always use **IQR**.

**Gradient Boosting (Recommended):**
- [Test Link: 65m² in Vračar (Sigma)](https://bg-real-estate.duckdns.org/predict/boost?sqm=65&rooms=2.5&floor=3&district=Vra%C4%8Dar&outlier_method=sigma)
- [Test Link: 65m² in Vračar (IQR)](https://bg-real-estate.duckdns.org/predict/boost?sqm=65&rooms=2.5&floor=3&district=Vra%C4%8Dar&outlier_method=iqr)

**K-Nearest Neighbors:**
- [Test Link: 40m² in Zemun](https://bg-real-estate.duckdns.org/predict/knn?sqm=40&rooms=1&district=Zemun&exclude_outliers=true&outlier_method=iqr)

**Decision Tree:**
- [Test Link: 120m² in Savski Venac](https://bg-real-estate.duckdns.org/predict/tree?sqm=120&rooms=4&district=Savski+venac)

**Polynomial Regression:**
- [Test Link: 55m² in Zvezdara](https://bg-real-estate.duckdns.org/predict?sqm=55&district=Zvezdara&round=2)

### Analytics and Statistics

**Advanced Analysis & Market Trends:**
Returns normality checks, distribution data, and the **monthly market trend**.
- [Test Link: Trends & Stats for Palilula](https://bg-real-estate.duckdns.org/analyze?fields=price&district=Palilula&outlier_method=iqr)

**Full Dataset Export:**
Returns statistics and correlation matrix for the district in one request.
- [Test Link: Full Data Dump (Zemun)](https://bg-real-estate.duckdns.org/full?district=Zemun)

**Feature Correlation:**
Matrix of how price, area, and rooms relate to each other.
- [Test Link: Correlation in Voždovac](https://bg-real-estate.duckdns.org/correlation?district=Vo%C5%BEdovac)

**Field Statistics:**
Basic stats (min, max, mean, median) for a specific metric.
- [Test Link: Statistics for Čukarica](https://bg-real-estate.duckdns.org/stats?district=%C4%8Cukarica&exclude_outliers=true)
