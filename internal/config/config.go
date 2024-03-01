package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

const configPath string = "./config/config.yaml"

type Config struct {
	StoragePath string `yaml:"storage_path"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func NewConfig() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read conifg: %s", err)
	}
	return &cfg
}
