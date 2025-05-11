package handlers

import "gorm.io/gorm"

type Handlers struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Handlers {
	return &Handlers{DB: db}
}
