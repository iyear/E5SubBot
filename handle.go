package main

import (
	"fmt"
	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

const (
	bStartContent string = "欢迎使用E5SubBot!\n请输入命令以启用"
)

var (
	UserStatus map[int64]int
	BindMaxNum int
)

const (
	USNone = iota
	USUnbind
	USWillBind
	USBind
)

func init() {
	//read config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	CheckErr(err)

	BindMaxNum = viper.GetInt("bindmax")

	UserStatus = make(map[int64]int)
}
func bStart(m *tb.Message) {
	bot.Send(m.Sender, bStartContent)
}
func bMy(m *tb.Message) {
	data := QueryData(db, m.Chat.ID)
	var inlineKeys = [][]tb.InlineButton{}
	for _, u := range data {
		inlineBtn := tb.InlineButton{
			Unique: u.refreshToken,
			Text:   u.refreshToken,
			Data:   u.uptime.Format("2006-01-02 15:04:05"),
		}
		bot.Handle(&inlineBtn, bMyinlineBtn)
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}
	bot.Send(m.Chat, "Hello!", &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
}
func bMyinlineBtn(c *tb.Callback) {
	bot.Send(c.Message.Chat, c.Data)
	bot.Respond(c)
}
func bBind(m *tb.Message) {
	tgId := m.Chat.ID
	fmt.Println("Auth: " + strconv.FormatInt(tgId, 10))
	bot.Send(m.Chat, "授权链接： [点击直达](https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=4d7c8a8a-0baf-497e-9608-57d6abfccce7&response_type=code&redirect_uri=http%3A%2F%2Flocalhost%2Fe5sub%2F&response_mode=query&scope=openid%20offline_access%20mail.read%20user.read)", tb.ModeMarkdown)
	_, err := bot.Send(m.Chat, "授权后回复整个http://localhost", &tb.ReplyMarkup{ForceReply: true})
	if err == nil {
		UserStatus[m.Chat.ID] = USWillBind
	}

}
func bAbout(m *tb.Message) {
	bot.Send(m.Sender, bStartContent)
}
func bOnText(m *tb.Message) {
	switch UserStatus[m.Chat.ID] {
	case USNone:
		{
			bot.Send(m.Chat, "发送/bind开始绑定嗷")
			return
		}
	case USWillBind:
		{
			if GetBindNum(m.Chat.ID) == BindMaxNum {
				bot.Send(m.Chat, "已经达到最大可绑定数")
				return
			}
			bot.Send(m.Chat, "正在绑定中……")
			info := BindUser(m)
			if info == "" {
				bot.Send(m.Chat, "绑定成功!")
			} else {
				bot.Send(m.Chat, info)
			}
		}
	case USBind:

	}
}
