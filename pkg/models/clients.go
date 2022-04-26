package models

type Client struct {
	ID           int64  `gorm:"primaryKey;autoIncrement;not null;index"`
	TgId         int64  `gorm:"not null;index"`
	RefreshToken string `gorm:"not null"`
	MsId         string `gorm:"not null"`
	Uptime       int64  `gorm:"autoUpdateTime;not null;index"`
	Alias        string `gorm:"not null"`
	ClientId     string `gorm:"not null"`
	ClientSecret string `gorm:"not null"`
	Other        string
}
