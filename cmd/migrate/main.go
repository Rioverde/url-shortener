package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/Rioverde/url-shortener/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// migrationsPath is where the SQL files live, relative to the project root.
// We pass it to migrate as a file:// URL.
const migrationsPath = "file://migrations"

func main() {
	// Parse the -cmd flag. Accepted values: up, down, version.
	cmd := flag.String("cmd", "up", "migration command: up | down | version")
	flag.Parse()

	// Load the same config the server uses, so the DB path is a single source of truth.
	cfg := config.MustLoad()

	// Build the database URL for golang-migrate. The sqlite3 driver expects sqlite3://<path>.
	dbURL := fmt.Sprintf("sqlite3://%s", cfg.StoragePath)

	// Create a migrate instance bound to our migrations folder and target DB.
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("failed to init migrator: %v", err)
	}
	// Always close — it releases the DB lock and closes the connection.
	defer m.Close()

	// Dispatch on the command we received.
	switch *cmd {
	case "up":
		// Apply every pending migration. ErrNoChange means we're already up to date.
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate up failed: %v", err)
		}
		fmt.Println("migrations applied")

	case "down":
		// Revert exactly one migration. Use Down() (without arg) to revert all.
		if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate down failed: %v", err)
		}
		fmt.Println("one migration reverted")

	case "version":
		// Print current schema version and dirty flag.
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("failed to get version: %v", err)
		}
		fmt.Printf("version=%d dirty=%v\n", version, dirty)

	default:
		log.Fatalf("unknown command %q (use up | down | version)", *cmd)
	}
}
