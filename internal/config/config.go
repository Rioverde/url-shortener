package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string           `yaml:"env" env-default:"local"`
	StoragePath string           `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServerConfig `yaml:"http_server"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	// Set the config path from the environment variable, or use a default if not set
	path := os.Getenv("CONFIG_PATH")
	// If CONFIG_PATH is not set, use the default path
	if path == "" {
		// TODO: Change to the main logger when it's implemented
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", path)
	}

	// Create a variable to hold the config
	var cnf Config

	// Load the config from the specified path
	if err := cleanenv.ReadConfig(path, &cnf); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	return &cnf
}
