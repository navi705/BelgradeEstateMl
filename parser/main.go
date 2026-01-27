package main

import (
	"io"
	"log/slog"
	"os"

	"sync"
	"time"
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
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 99; i++ {
			estates, err := FourZidaList(i)
			if err != nil {
				slog.Error("Error parsing 4zida", "page", i, "error", err)
				continue
			}
			for _, e := range estates {
				if err := s.SaveEstate(e); err != nil {
					slog.Error("Error saving estate from 4zida", "link", e.Link, "error", err)
				}
			}
			slog.Info("Saved page from 4zida", "page", i, "count", len(estates))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		page := 1
		for {
			estates, err := HaloOglasiList(page)
			if err != nil {
				slog.Error("Error parsing HaloOglasi", "page", page, "error", err)
				break
			}
			if len(estates) == 0 {
				break
			}
			for _, e := range estates {
				if err := s.SaveEstate(e); err != nil {
					slog.Error("Error saving estate from HaloOglasi", "link", e.Link, "error", err)
				}
			}
			slog.Info("Saved page from HaloOglasi", "page", page, "count", len(estates))
			page++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		page := 1
		for {
			estates, err := NekretnineList(page)
			if err != nil {
				slog.Error("Error parsing Nekretnine", "page", page, "error", err)
				break
			}
			if len(estates) == 0 {
				break
			}
			for _, e := range estates {
				if err := s.SaveEstate(e); err != nil {
					slog.Error("Error saving estate from Nekretnine", "link", e.Link, "error", err)
				}
			}
			slog.Info("Saved page from Nekretnine", "page", page, "count", len(estates))
			page++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		page := 1
		for {
			estates, err := CityExpertList(page)
			if err != nil {
				slog.Error("Error parsing CityExpert", "page", page, "error", err)
				break
			}
			if len(estates) == 0 {
				break
			}
			for _, e := range estates {
				if err := s.SaveEstate(e); err != nil {
					slog.Error("Error saving estate from CityExpert", "link", e.Link, "error", err)
				}
			}
			slog.Info("Saved page from CityExpert", "page", page, "count", len(estates))
			page++
		}
	}()

	wg.Wait()
	slog.Info("Parser run completed")
}

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		slog.Error("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	storage, err := NewStorage(connStr)
	if err != nil {
		slog.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}

	if err := storage.Migrate(); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	runParser(storage)

	ticker := time.NewTicker(48 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		runParser(storage)
	}
}
