package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const maxLogSize = 40 * 1024 * 1024 // 40MB

type LogRotator struct {
	Filename string
	MaxSize  int64
	file     *os.File
	mu       sync.Mutex
}

func (l *LogRotator) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		f, err := os.OpenFile(l.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return 0, err
		}
		l.file = f
	}

	fi, err := l.file.Stat()
	if err == nil && fi.Size() > l.MaxSize {
		l.file.Close()
		f, err := os.OpenFile(l.Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			return 0, err
		}
		l.file = f
	}

	return l.file.Write(p)
}

func init() {
	logRotator := &LogRotator{
		Filename: "parser.log",
		MaxSize:  maxLogSize,
	}

	multiWriter := io.MultiWriter(os.Stdout, logRotator)
	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(handler))

	slog.Info("parser initialized")
}

func runParser(s *Storage) {
	slog.Info("Starting parser run...")

	sites := []struct {
		name string
		fn   func(int) ([]RealEstate, error)
		max  int // 0 means until no more elements
	}{
		{"4zida.rs", FourZidaList, 99},
		{"halooglasi.com", HaloOglasiList, 0},
		{"nekretnine.rs", NekretnineList, 0},
		{"cityexpert.rs", CityExpertList, 0},
	}

	var wg sync.WaitGroup
	for _, site := range sites {
		wg.Add(1)
		go func(sName string, sFn func(int) ([]RealEstate, error), sMax int) {
			defer wg.Done()
			start := time.Now()

			page := 1
			for {
				if sMax > 0 && page > sMax {
					break
				}

				estates, err := sFn(page)
				if err != nil {
					slog.Error("Error parsing", "site", sName, "page", page, "error", err)
					parserErrors.WithLabelValues(sName, "list_fetch").Inc()
					break
				}

				if len(estates) == 0 {
					break
				}

				for _, e := range estates {
					if err := s.SaveEstate(e); err != nil {
						slog.Error("Error saving estate", "site", sName, "link", e.Link, "error", err)
						parserErrors.WithLabelValues(sName, "db_save").Inc()
					} else {
						processedItems.WithLabelValues(sName, "processed").Inc()
					}
				}

				slog.Info("Saved page", "site", sName, "page", page, "count", len(estates))
				page++
			}

			duration := time.Since(start).Seconds()
			runDuration.WithLabelValues(sName).Observe(duration)
			lastRunTimestamp.WithLabelValues(sName).SetToCurrentTime()
			slog.Info("Site parsing completed", "site", sName, "duration", duration)
		}(site.name, site.fn, site.max)
	}

	wg.Wait()
	slog.Info("Parser run completed")
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("PROJECT_USER")
	dbPass := os.Getenv("PROJECT_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSL := os.Getenv("DB_SSL_MODE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" {
		slog.Error("Database environment variables are not fully set")
		os.Exit(1)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, dbSSL)

	storage, err := NewStorage(connStr)
	if err != nil {
		slog.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}

	if err := storage.Migrate(); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Start Prometheus metrics server
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "2112" // fallback
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		slog.Info("Starting metrics server", "port", metricsPort)
		if err := http.ListenAndServe(":"+metricsPort, nil); err != nil {
			slog.Error("Metrics server failed", "error", err)
		}
	}()

	runParser(storage)

	ticker := time.NewTicker(48 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		runParser(storage)
	}
}
