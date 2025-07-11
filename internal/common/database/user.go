package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramID  int64
	Channels    []*Channel `gorm:"many2many:user_channel;constraint:OnDelete:CASCADE;"`
	Guilds      []*Guild   `gorm:"many2many:user_guild;constraint:OnDelete:CASCADE;"`
	LastGuildID string
}
