# 🏙 Belgrade Estate ML Proxy

High-performance caching proxy for the Belgrade Estate ML engine.

## Features
- **Parallel Processing**: Built on Go's high-performance standard library; every request is handled in a separate goroutine.
- **LRU Caching**: Stores the last 20 unique prediction results in memory.
- **Anti-Spam (Rate Limiting)**: 
    - **Limit**: 2 requests per second per IP.
    - **Burst**: Allows up to 5 requests at once before blocking.
    - **Action**: Returns `429 Too Many Requests`.
- **Advanced Monitoring**: Dedicated metrics for unique IPs and traffic patterns.

## 🚀 Running with Docker

```bash
docker run -p 8000:8000 \
  -e ML_ENGINE_URL="http://ml-engine:8080" \
  ghcr.io/navi705/belgrade-estate-server:latest
```

## 📊 Monitoring
Metrics are exposed at `/metrics`. 
See `grafana/dashboard.json` for the professional dashboard template.
