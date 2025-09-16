package handlers

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type handlers struct {
	db     *gorm.DB
	logger *zap.Logger
}

func New(db *gorm.DB, logger *zap.Logger) *handlers {
	return &handlers{db, logger.Named("handlers")}
}
