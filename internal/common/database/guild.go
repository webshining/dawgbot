package database

type Guild struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	Channels []*Channel
}
