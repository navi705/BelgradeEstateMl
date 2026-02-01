package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	processedItems = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "parser_items_processed_total",
		Help: "Total number of real estate items processed",
	}, []string{"site", "status"})

	parserErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "parser_errors_total",
		Help: "Total number of errors during parsing",
	}, []string{"site", "phase"})

	lastRunTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "parser_last_run_timestamp_seconds",
		Help: "Unix timestamp of the last successful run per site",
	}, []string{"site"})

	runDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "parser_run_duration_seconds",
		Help:    "Duration of the parser run in seconds",
		Buckets: []float64{10, 30, 60, 120, 300, 600, 1200, 3600},
	}, []string{"site"})
)
