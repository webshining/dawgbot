package middlewares

import "gorm.io/gorm"

type Middlewares struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Middlewares {
	return &Middlewares{DB: db}
}
