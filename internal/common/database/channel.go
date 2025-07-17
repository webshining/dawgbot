package database

type Channel struct {
	ID      string `gorm:"primaryKey"`
	Name    string
	GuildID string
	Users   []*User `gorm:"many2many:user_channel;constraint:OnDelete:CASCADE;"`
}
