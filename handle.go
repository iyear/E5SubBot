package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

const (
	bStartContent string = "欢迎使用E5SubBot!\n请输入命令以启用Bot"
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
	data := QueryDataByTG(db, m.Chat.ID)
	var inlineKeys [][]tb.InlineButton
	for _, u := range data {
		inlineBtn := tb.InlineButton{
			Unique: "my" + u.msId,
			Text:   gjson.Get(u.other, "alias").String(),
			Data:   u.msId,
		}
		bot.Handle(&inlineBtn, bMyInlineBtn)
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}
	_, err := bot.Send(m.Chat, "选择一个账户查看具体信息\n\n绑定数: "+strconv.Itoa(GetBindNum(m.Chat.ID))+"/"+strconv.Itoa(BindMaxNum), &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
	fmt.Println(err)
}
func bMyInlineBtn(c *tb.Callback) {
	fmt.Println(c.Data)
	r := QueryDataByMS(db, c.Data)
	u := r[0]
	bot.Send(c.Message.Chat, "信息\n别名："+gjson.Get(u.other, "alias").String()+"\nMS_ID(MD5): "+u.msId+"\n最近更新时间: "+u.uptime.Format("2006-01-02 15:04:05"))
	bot.Respond(c)
}
func bBind(m *tb.Message) {
	tgId := m.Chat.ID
	fmt.Println("Auth: " + strconv.FormatInt(tgId, 10))
	bot.Send(m.Chat, "授权链接： [点击直达]("+authUrl+")", tb.ModeMarkdown)
	_, err := bot.Send(m.Chat, "回复格式：http://localhost/...+空格+别名(用于管理)", &tb.ReplyMarkup{ForceReply: true})
	if err == nil {
		UserStatus[m.Chat.ID] = USWillBind
	}

}
func bUnBind(m *tb.Message) {
	data := QueryDataByTG(db, m.Chat.ID)
	var inlineKeys [][]tb.InlineButton
	for _, u := range data {
		inlineBtn := tb.InlineButton{
			Unique: "unbind" + u.msId,
			Text:   gjson.Get(u.other, "alias").String(),
			Data:   u.msId,
		}
		bot.Handle(&inlineBtn, bUnBindInlineBtn)
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}
	bot.Send(m.Chat, "选择一个账户将其解绑\n\n当前绑定数: "+strconv.Itoa(GetBindNum(m.Chat.ID))+"/"+strconv.Itoa(BindMaxNum), &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
}
func bUnBindInlineBtn(c *tb.Callback) {
	fmt.Println(c.Data)
	r := QueryDataByMS(db, c.Data)
	u := r[0]
	if ok, _ := DelData(db, u.msId); !ok {
		fmt.Println(u.msId + " UnBind ERROR")
		bot.Send(c.Message.Chat, "解绑失败!")
		return
	}
	fmt.Println(u.msId + " UnBind Success")
	bot.Send(c.Message.Chat, "解绑成功!")
	bot.Respond(c)
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
			UserStatus[m.Chat.ID] = USNone
		}
	}
}
