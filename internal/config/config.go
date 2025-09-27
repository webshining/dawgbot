package config

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

var (
	k        = koanf.New(".")
	validate = validator.New()
)

func Load() (*Config, error) {
	var cfg Config

	k.Load(confmap.Provider(defaultConfig, "_"), nil)
	k.Load(file.Provider(".env"), dotenv.ParserEnv("", "_", func(k string) string {
		return k
	}))
	err := errors.Join(
		k.Unmarshal("", &cfg),
		validate.Struct(&cfg),
	)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func MustLoad(logger *zap.Logger) *Config {
	cfg, err := Load()
	if err != nil {
		logger.Named("Config").Fatal("failed to load config", zap.Error(err))
	}

	return cfg
}
