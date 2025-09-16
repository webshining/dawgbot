package middlewares

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type middlewares struct {
	db     *gorm.DB
	logger *zap.Logger
}

func New(db *gorm.DB, logger *zap.Logger) *middlewares {
	return &middlewares{db, logger.Named("middlewares")}
}
