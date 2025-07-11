//go:build !prod

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("database.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Guild{}, &Channel{}, &User{})
	return db, nil
}
