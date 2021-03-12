package core

type Client struct {
	TgId         int64  `gorm:"column:tg_id"`
	RefreshToken string `gorm:"column:refresh_token"`
	MsId         string `gorm:"column:ms_id"`
	Uptime       int64  `gorm:"column:uptime"`
	Alias        string `gorm:"column:alias"`
	ClientId     string `gorm:"column:client_id"`
	ClientSecret string `gorm:"column:client_secret"`
	Other        string `gorm:"column:other"`
}
