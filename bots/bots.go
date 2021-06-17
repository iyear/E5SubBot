package bots

import (
	"fmt"
	"github.com/iyear/E5SubBot/config"
	"github.com/iyear/E5SubBot/logger"
	"github.com/iyear/E5SubBot/model"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
	"strconv"
	"time"
)

var (
	bot *tb.Bot
)

const (
	logo = `
  ______ _____ _____       _     ____        _   
 |  ____| ____/ ____|     | |   |  _ \      | |  
 | |__  | |__| (___  _   _| |__ | |_) | ___ | |_ 
 |  __| |___ \\___ \| | | | '_ \|  _ < / _ \| __|
 | |____ ___) |___) | |_| | |_) | |_) | (_) | |_ 
 |______|____/_____/ \__,_|_.__/|____/ \___/ \__|
`
)

func BotStart() {
	var err error
	fmt.Println(logo)
	//read config
	config.InitConfig()
	//Init Logger
	logger.InitLogger()
	//InitDB
	model.InitDB()
	//Init Task
	InitTask()

	Poller := &tb.LongPoller{Timeout: 15 * time.Second}
	spamPoller := tb.NewMiddlewarePoller(Poller, func(upd *tb.Update) bool {
		if upd.Message == nil {
			return true
		}
		if !upd.Message.Private() {
			return false
		}
		return true
	})
	botSetting := tb.Settings{
		Token:  config.BotToken,
		Poller: spamPoller,
	}
	//set socks5
	if config.Socks5 != "" {
		fmt.Println("Proxy:" + config.Socks5)
		dialer, err := proxy.SOCKS5("tcp", config.Socks5, nil, proxy.Direct)
		if err != nil {
			zap.S().Errorw("failed to make dialer", "error", err, "socks5", config.Socks5)
		}
		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		httpTransport.Dial = dialer.Dial
		botSetting.Client = httpClient
	}
	//create bot
	bot, err = tb.NewBot(botSetting)
	if err != nil {
		zap.S().Errorw("failed to create bot", "error", err)
		return
	}
	fmt.Println("Bot: " + strconv.Itoa(bot.Me.ID) + " " + bot.Me.Username)

	MakeHandle()
	fmt.Println("Bot Start")
	fmt.Println("------------")
	bot.Start()
}
func MakeHandle() {
	//所有用户
	bot.Handle("/start", bStart)
	bot.Handle("/my", bMy)
	bot.Handle("/bind", bBind1)
	bot.Handle("/unbind", bUnBind)
	bot.Handle("/export", bExport)
	bot.Handle("/help", bHelp)
	bot.Handle(tb.OnText, bOnText)
	//管理员
	bot.Handle("/task", bTask)
	bot.Handle("/log", bLog)
}
