package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mlURL *url.URL
	cache *LRUCache
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/metrics" {
		promhttp.Handler().ServeHTTP(w, r)
		return
	}

	cacheKey := r.URL.String()
	if val, ok := cache.Get(cacheKey); ok {
		metricCacheHits.Inc()
		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("Content-Type", "application/json")
		w.Write(val)
		return
	}

	// Create proxy request
	proxyURL := *mlURL
	proxyURL.Path = r.URL.Path
	proxyURL.RawQuery = r.URL.RawQuery

	req, _ := http.NewRequest(r.Method, proxyURL.String(), r.Body)
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "ML Engine unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Cache successful responses
	if resp.StatusCode == http.StatusOK {
		cache.Put(cacheKey, body)
	}

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	metricTotalRequests.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", resp.StatusCode)).Inc()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mURL := os.Getenv("ML_ENGINE_URL")
	if mURL == "" {
		mURL = "http://localhost:8080"
	}

	var err error
	mlURL, err = url.Parse(mURL)
	if err != nil {
		log.Fatal(err)
	}

	cache = NewLRUCache(20)
	stats := NewIPStats()

	mux := http.NewServeMux()
	mux.HandleFunc("/", proxyHandler)

	handler := stats.LimitMiddleware(mux)

	log.Printf("Proxy server starting on :%s, forwarding to %s", port, mURL)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
