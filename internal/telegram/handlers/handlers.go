package handlers

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handlers struct {
	DB     *gorm.DB
	logger *zap.Logger
}

func New(db *gorm.DB, logger *zap.Logger) *Handlers {
	return &Handlers{DB: db, logger: logger}
}
