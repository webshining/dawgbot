package database

type User struct {
	ID          int64      `gorm:"primaryKey"`
	Guilds      []*Guild   `gorm:"many2many:user_guild;constraint:OnDelete:CASCADE;"`
	Channels    []*Channel `gorm:"many2many:user_channel;constraint:OnDelete:CASCADE;"`
	LastGuildID string
}
