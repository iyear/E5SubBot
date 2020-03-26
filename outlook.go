package main

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	cliid   string
	rediuri string
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	cliid = viper.GetString("client_id")
	rediuri = viper.GetString("redirect_uri")
}
func GetToken(code string) string {
	return ""
}
