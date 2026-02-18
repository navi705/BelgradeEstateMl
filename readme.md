# Belgrade Estate ML

Professional real estate analytics and predictive engine for Belgrade.

## 🏗 Architecture
1.  **Parser**: Collects data from Belgrade real estate sites.
2.  **Database (PostgreSQL)**: Stores cleaned and standardized data.
3.  **ML Engine (ml/)**: Core mathematical models and analytics.
4.  **Proxy Gateway (server/)**: High-performance "fat" proxy with:
    - **LRU Caching**: Instant response for repeated queries.
    - **Rate Limiting**: IP-based anti-spam protection.
    - **Traffic Monitoring**: Tracking unique users and cache hits.

## 🚀 Public API
The gateway is available on port `8000` (by default). For direct ML engine access, use port `8080`.

---
