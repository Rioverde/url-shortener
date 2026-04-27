package server

import (
	"crypto/rand"
	"log/slog"
	"os"

	"github.com/Rioverde/url-shortener/internal/config"
	"github.com/Rioverde/url-shortener/internal/service"
	"github.com/Rioverde/url-shortener/internal/storage"
)

const (
	prodEnv = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting URL Shortener Server", "env", cfg.Env)

	// TODO(Init Sorage): init Storage

	// TODO(Init Sorage): init Logger

	repo := storage.NewMapStorage()

	gen := service.NewCryptoGenerator(rand.Reader)

	svc := service.NewURLService(repo, gen)

	_ = svc

	// TODO(Init Sorage): init Server

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
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
