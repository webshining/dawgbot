package database

import "gorm.io/gorm"

type Guild struct {
	gorm.Model
	ID        string
	Name      string
	Channels  []*Channel  `gorm:"foreignKey:GuildID;constraint:OnDelete:CASCADE;"`
	Users     []*User     `gorm:"many2many:user_guild;constraint:OnDelete:CASCADE;"`
	Playlists []*Playlist `gorm:"foreignKey:GuildID;constraint:OnDelete:CASCADE;"`
}
