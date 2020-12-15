package config

import (
	"fmt"
	"os"

	logger "github.com/moguchev/service/pkg/logger"
	"github.com/moguchev/service/pkg/pgsql"
	"gopkg.in/yaml.v2"
)

type (
	ServerConfig struct {
		Address     string
		APIBasePath string
	}

	Config struct {
		Server *ServerConfig  `yaml:"server"`
		DB     *pgsql.Config  `yaml:"db"`
		Log    *logger.Config `yaml:"log"`
	}
)

func GetConfig(path string) (Config, error) {
	var cfg Config

	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("open config file: %w", err)
	}

	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("decode: %w", err)
	}

	return cfg, nil
}
