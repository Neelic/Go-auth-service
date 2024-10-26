package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Println("CONFIG_PATH is empty")
		b, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		configPath = filepath.Join(b, "config/local.yaml")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("config file error: %s", err)
	}

	return &cfg
}
