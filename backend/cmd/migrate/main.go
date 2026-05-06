package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	var migrationDir string
	flag.StringVar(&migrationDir, "dir", "internal/infrastructure/persistence/migrations", "directory for migration files")
	flag.Parse()

	// Build DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	m, err := migrate.New(
		"file://"+migrationDir,
		dsn,
	)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	// Read command from args
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Please specify a command: 'up' or 'down'")
	}

	command := args[0]
	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration UP successful!")
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration DOWN successful!")
	case "force":
		if len(args) < 2 {
			log.Fatal("Please specify a version to force")
		}
		var version int
		fmt.Sscanf(args[1], "%d", &version)
		if err := m.Force(version); err != nil {
			log.Fatalf("Force failed: %v", err)
		}
		log.Printf("Forced to version %d", version)
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Could not get version: %v", err)
		}
		log.Printf("Current version: %d, Dirty: %v", v, dirty)
	default:
		log.Fatalf("Unknown command: %s. Use 'up', 'down', 'force', or 'version'", command)
	}
}
