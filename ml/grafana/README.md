## 🚀 Quick Start (Dashboard Import)

1.  **Open Grafana**: Go to **Dashboards** -> **Import**.
2.  **Upload JSON**: Select the `dashboard.json` file from this folder.
3.  **Configure Data Source**: You will see a dropdown asking for a **Prometheus** source. Select your Prometheus data source there.
4.  **Import**: Click the import button.

---

## 🔌 Data Sources Setup

For the best experience, add **two** data sources in Grafana:
1.  **Prometheus**: URL `http://localhost:8080/metrics` (or your server IP).
2.  **PostgreSQL**: Your database connection details.

---

## 🛠 Variables for Interactivity
The dashboard uses a `$district` variable. It is populated automatically from the metrics:
- **Type**: Query
- **Data Source**: Prometheus
- **Query**: `label_values(realestate_api_requests_total, district)`

---

## 📝 Manual Queries (If creating custom panels)

### Prediction Card (Prometheus)
```promql
realestate_prediction_price_last{algorithm="boost", district="$district"}
```

### Scatter Plot (SQL)
```sql
SELECT
  square_meter as "x",
  price as "y",
  district as "label"
FROM estates
WHERE district = '${district}'
ORDER BY 1
```
