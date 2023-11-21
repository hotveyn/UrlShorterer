package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	Storage    Storage    `yaml:"storage" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server" env-required:"true"`
}

type Storage struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	StorageName string `yaml:"storage_name" env-required:"true"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	Port        int           `yaml:"port" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func (s *HTTPServer) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is required")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config by path %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}
