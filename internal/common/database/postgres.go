//go:build prod

package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Guild{}, &Channel{})
	return db, nil
}
