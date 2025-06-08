package database

import "gorm.io/gorm"

type Playlist struct {
	gorm.Model
	GuildID string
	FileUrl string
}
