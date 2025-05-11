package database

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	ID      string
	Name    string
	GuildID string
	Users   []*User `gorm:"many2many:user_channel;constraint:OnDelete:CASCADE;"`
}
