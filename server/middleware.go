package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/time/rate"
)

var (
	metricTotalRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "proxy_requests_total",
		Help: "The total number of requests handled by the proxy",
	}, []string{"method", "path", "status"})

	metricUniqueIPs = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proxy_unique_ips_count",
		Help: "Total number of unique IPs seen since startup",
	})

	metricCacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "proxy_cache_hits_total",
		Help: "Total number of cache hits",
	})
)

type IPStats struct {
	uniqueIPs map[string]bool
	limiter   map[string]*rate.Limiter
	mu        sync.Mutex
}

func NewIPStats() *IPStats {
	return &IPStats{
		uniqueIPs: make(map[string]bool),
		limiter:   make(map[string]*rate.Limiter),
	}
}

func (s *IPStats) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		// Strip port if exists
		for i := len(ip) - 1; i >= 0; i-- {
			if ip[i] == ':' {
				ip = ip[:i]
				break
			}
		}

		s.mu.Lock()
		if !s.uniqueIPs[ip] {
			s.uniqueIPs[ip] = true
			metricUniqueIPs.Set(float64(len(s.uniqueIPs)))
		}

		limiter, exists := s.limiter[ip]
		if !exists {
			// Limit to 2 requests per second with a burst of 5
			limiter = rate.NewLimiter(rate.Every(time.Second/2), 5)
			s.limiter[ip] = limiter
		}
		s.mu.Unlock()

		if !limiter.Allow() {
			log.Printf("Rate limit exceeded for IP: %s", ip)
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
