package model

import (
	"github.com/iyear/E5SubBot/config"
)

type Client struct {
	ID           int    `gorm:"unique;primaryKey;not null"`
	TgId         int64  `gorm:"not null"`
	RefreshToken string `gorm:"not null"`
	MsId         string `gorm:"not null"`
	Uptime       int64  `gorm:"autoUpdateTime;not null"`
	Alias        string `gorm:"not null"`
	ClientId     string `gorm:"not null"`
	ClientSecret string `gorm:"not null"`
	Other        string
}

func (c *Client) TableName() string {
	return config.Table
}
func NewClient(clientId string, clientSecret string) *Client {
	return &Client{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}
