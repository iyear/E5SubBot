package main

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	BotToken string
	Socks5   string
)

const (
	dbDriverName = "sqlite3"
	dbName       = "./data.db"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	CheckErr(err)
	BotToken = viper.GetString("bot_token")
	Socks5 = viper.GetString("socks5")

}
func main() {
	botsettings := tb.Settings{
		Token:  BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}
	if Socks5 != "" {
		fmt.Println("Proxy:" + Socks5)
		dialer, err := proxy.SOCKS5("tcp", Socks5, nil, proxy.Direct)
		if err != nil {
			log.Fatal("Error creating dialer, aborting.")
		}
		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		httpTransport.Dial = dialer.Dial
		botsettings.Client = httpClient
	}
	db, err := sql.Open(dbDriverName, dbName)
	CheckErr(err)
	if !FileExist(dbName) {
		CreateTB(db)
	}
	b, err := tb.NewBot(botsettings)
	CheckErr(err)
	//b.Handle(tb.OnText, func(m *tb.Message) {
	//	b.Send(m.Sender, "hello world")
	//})

	b.Start()
}
func CheckErr(err error) bool {
	if err != nil {
		log.Println(err)
		fmt.Println("error: ", err.Error())
		panic(err)
		return false
	}
	return true
}
func FileExist(Path string) bool {
	if _, err := os.Stat(Path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			CheckErr(err)
		}
	}
	return true
}
