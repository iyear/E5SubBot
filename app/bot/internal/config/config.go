package config

import (
	"github.com/spf13/viper"
)

var C struct {
	Bot struct {
		Token  string `mapstructure:"token"`
		Socks5 struct {
			Enable   bool   `mapstructure:"enable"`
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
		} `mapstructure:"socks5"`
		Admin []int64 `mapstructure:"admin"`
	} `mapstructure:"bot"`
	Biz struct {
		Notice string `mapstructure:"notice"`
	}
	Ctrl struct {
		BindMax     int    `mapstructure:"bind_max"`
		Goroutine   int    `mapstructure:"goroutine"`
		ErrorLimit  int    `mapstructure:"error_limit"`
		Cron        string `mapstructure:"cron"`
		DefaultLang string `mapstructure:"default_lang"`
	}
	DB struct {
		Driver string `mapstructure:"driver"`
		MySQL  struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			Database string `mapstructure:"database"`
		} `mapstructure:"mysql"`
	}
}

func Init(path string) error {
	c := viper.New()
	c.SetConfigFile(path)
	if err := c.ReadInConfig(); err != nil {
		return err
	}

	return c.Unmarshal(&C)
}
