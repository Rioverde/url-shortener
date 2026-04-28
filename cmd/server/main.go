package main

import (
	"crypto/rand"
	"log/slog"
	"os"

	mvLogger "github.com/Rioverde/url-shortener/internal/api/middleware"
	"github.com/Rioverde/url-shortener/internal/config"
	"github.com/Rioverde/url-shortener/internal/domain"
	"github.com/Rioverde/url-shortener/internal/lib/codegen"
	"github.com/Rioverde/url-shortener/internal/repo/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	prodEnv = "prod"
)

func main() {

	// Load the configuration
	cfg := config.MustLoad()

	// Set up the logger based on the environment
	log := setupLogger(cfg.Env)

	// Log the startup message with the environment
	log.Info("Starting URL Shortener Server", "env", cfg.Env)

	// Initialize the SQLite storage. If it fails, log the error and exit.
	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to create SQLite storage", "error", err)
		os.Exit(1)
	}

	// Ensure the database connection is closed when main exits. Log any error that occurs during close.
	defer func() {
		if err := storage.Close(); err != nil {
			log.Error("failed to close database connection", "error", err)
		}
	}()

	// Initialize the code generator using crypto/rand as the entropy source
	gen := codegen.NewCryptoGenerator(rand.Reader)

	// Initialize the URL service with the storage and code generator
	svc := domain.NewURLService(storage, gen)

	_ = svc

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mvLogger.New(log))
	r.Use(middleware.Recoverer)

}

// setupLogger sets up the logger based on the environment
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	// Set the logger based on the environment
	switch env {
	case prodEnv:
		// In production, we want a structured logger with JSON output
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		// In local and dev environments, we want a human-readable logger with debug level
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
