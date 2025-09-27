package database

import (
	"bot/internal/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(config config.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch config.Driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.Url))
	default:
		db, err = gorm.Open(sqlite.Open(config.Url))
	}
	db.AutoMigrate(&Channel{}, &Guild{}, &User{})

	return db, err
}

func MustConnect(config config.DatabaseConfig, logger *zap.Logger) *gorm.DB {
	db, err := Connect(config)
	if err != nil {
		logger.Named("database").Fatal("failed to connect to database", zap.String("driver", config.Driver), zap.String("url", config.Url), zap.Error(err))
	}

	return db
}
