package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env      string     `yaml:"env" env-default:"local"`
	DataBase Storage    `yaml:"database" env-required:"true"`
	Http     HTTPServer `yaml:"http_server" env-required:"true"`
}

type Storage struct {
	StorageURL     string `yaml:"storage_url" env-required:"true"`
	MigrationsPath string `yaml:"migrations_path" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
