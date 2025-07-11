//go:build prod

package database

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dns string) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	timeout := time.After(1 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
		if err == nil {
			db.AutoMigrate(&Guild{}, &Channel{}, &User{})
			return db, nil
		}

		select {
		case <-timeout:
			return nil, err
		case <-ticker.C:
			// next try
		}
	}

}
