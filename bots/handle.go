package bots

import (
	"encoding/json"
	"fmt"
	"github.com/iyear/E5SubBot/config"
	"github.com/iyear/E5SubBot/model"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	UserStatus       map[int64]int
	UserClientId     map[int64]string
	UserClientSecret map[int64]string
)

const (
	USNone = iota
	USBind1
	USBind2
)

func init() {
	UserStatus = make(map[int64]int)
	UserClientId = make(map[int64]string)
	UserClientSecret = make(map[int64]string)
}

func bStart(m *tb.Message) {
	bot.Send(m.Sender, config.WelcomeContent)
	bHelp(m)
}

func bMy(m *tb.Message) {
	var data []*model.Client
	model.DB.Where("tg_id = ?", m.Chat.ID).Find(&data)
	var inlineKeys [][]tb.InlineButton
	for _, u := range data {
		inlineBtn := tb.InlineButton{
			Unique: "my" + strconv.Itoa(u.ID),
			Text:   u.Alias,
			Data:   strconv.Itoa(u.ID),
		}
		bot.Handle(&inlineBtn, bMyInlineBtn)
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}

	bot.Send(m.Chat,
		fmt.Sprintf("✨ 选择一个账户查看具体信息\n\n绑定数: %d/%d", GetBindNum(m.Chat.ID), config.BindMaxNum),
		&tb.ReplyMarkup{InlineKeyboard: inlineKeys})
}
func bMyInlineBtn(c *tb.Callback) {
	var u *model.Client
	model.DB.Where("id = ?", c.Data).First(&u)
	bot.Send(c.Message.Chat,
		fmt.Sprintf("ℹ 信息\n别名：%s\nms_id: %s\nclient_id: %s\nclient_secret: %s\n最近更新时间: %s",
			u.Alias,
			u.MsId,
			u.ClientId,
			u.ClientSecret,
			time.Unix(u.Uptime, 0).Format("2006-01-02 15:04:05"),
		),
	)
	bot.Respond(c)
}

func bBind1(m *tb.Message) {
	bot.Send(m.Chat,
		fmt.Sprintf("👉 应用注册： [点击直达](%s)", model.GetMSRegisterAppUrl()),
		tb.ModeMarkdown,
	)

	bot.Send(m.Chat,
		"⚠ 请回复 `client_id(空格)client_secret`",
		&tb.SendOptions{ParseMode: tb.ModeMarkdown,
			ReplyMarkup: &tb.ReplyMarkup{ForceReply: true}},
	)

	UserStatus[m.Chat.ID] = USBind1
	UserClientId[m.Chat.ID] = m.Text
}
func bBind2(m *tb.Message) {
	tmp := strings.Split(m.Text, " ")
	if len(tmp) != 2 {
		bot.Send(m.Chat, "⚠ 错误的格式")
		return
	}
	ClientId := tmp[0]
	ClientSecret := tmp[1]
	bot.Send(m.Chat,
		"👉 授权账户： [点击直达]("+model.GetMSAuthUrl(ClientId)+")",
		tb.ModeMarkdown,
	)

	bot.Send(m.Chat,
		"⚠ 请回复`http://localhost/……(空格)别名`(用于管理)",
		&tb.SendOptions{ParseMode: tb.ModeMarkdown,
			ReplyMarkup: &tb.ReplyMarkup{ForceReply: true},
		},
	)

	UserStatus[m.Chat.ID] = USBind2
	UserClientId[m.Chat.ID] = ClientId
	UserClientSecret[m.Chat.ID] = ClientSecret
}

func bUnBind(m *tb.Message) {
	var data []*model.Client
	model.DB.Where("tg_id = ?", m.Chat.ID).Find(&data)
	var inlineKeys [][]tb.InlineButton

	for _, u := range data {
		inlineBtn := tb.InlineButton{
			Unique: "unbind" + strconv.Itoa(u.ID),
			Text:   u.Alias,
			Data:   strconv.Itoa(u.ID),
		}
		bot.Handle(&inlineBtn, bUnBindInlineBtn)
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}

	bot.Send(m.Chat,
		fmt.Sprintf("⚠ 选择一个账户将其解绑\n\n当前绑定数: %d/%d", GetBindNum(m.Chat.ID), config.BindMaxNum),
		&tb.ReplyMarkup{InlineKeyboard: inlineKeys},
	)
}
func bUnBindInlineBtn(c *tb.Callback) {
	if result := model.DB.Where("id = ?", c.Data).Delete(&model.Client{}); result.Error != nil {
		zap.S().Errorw("failed to delete db data",
			"error", result.Error,
			"id", c.Data,
		)
		bot.Send(c.Message.Chat, "⚠ 解绑失败!")
		return
	}
	bot.Send(c.Message.Chat, "✨ 解绑成功!")
	bot.Respond(c)
}
func bExport(m *tb.Message) {
	type ClientExport struct {
		Alias        string
		ClientId     string
		ClientSecret string
		RefreshToken string
		Other        string
	}
	var exports []ClientExport
	var data []*model.Client
	model.DB.Where("tg_id = ?", m.Chat.ID).Find(&data)
	if len(data) == 0 {
		bot.Send(m.Chat, "⚠ 你还没有绑定过账户嗷~")
		return
	}
	for _, u := range data {
		var cExport = ClientExport{
			Alias:        u.Alias,
			ClientId:     u.ClientId,
			ClientSecret: u.ClientSecret,
			RefreshToken: u.RefreshToken,
			Other:        u.Other,
		}
		exports = append(exports, cExport)
	}
	export, err := json.MarshalIndent(exports, "", "\t")
	if err != nil {
		zap.S().Errorw("failed to marshal json",
			"error", err)
		bot.Send(m.Chat, fmt.Sprintf("⚠ 获取JSON失败!\n\nERROR: %s", err.Error()))
		return
	}
	fileName := fmt.Sprintf("./%d_export_tmp.json", m.Chat.ID)
	if err = ioutil.WriteFile(fileName, export, 0644); err != nil {
		zap.S().Errorw("failed to write file",
			"error", err)
		bot.Send(m.Chat, "⚠ 写入临时文件失败~\n"+err.Error())
		return
	}
	exportFile := &tb.Document{
		File:     tb.FromDisk(fileName),
		FileName: strconv.FormatInt(m.Chat.ID, 10) + ".json",
		MIME:     "text/plain",
	}
	bot.Send(m.Chat, exportFile)

	if exportFile.InCloud() != true || os.Remove(fileName) != nil {
		zap.S().Errorw("failed to export files")
	}
}
func bHelp(m *tb.Message) {
	bot.Send(
		m.Sender,
		config.HelpContent+"\n"+config.Notice,
		&tb.SendOptions{DisableWebPagePreview: false})
}
func bOnText(m *tb.Message) {
	switch UserStatus[m.Chat.ID] {
	case USNone:
		{
			bot.Send(m.Chat, "✨ 发送 /help 获取帮助嗷")
			return
		}
	case USBind1:
		{
			if !m.IsReply() {
				bot.Send(m.Chat, "⚠ 请通过回复方式绑定")
				return
			}
			bBind2(m)
		}
	case USBind2:
		{
			if !m.IsReply() {
				bot.Send(m.Chat, "⚠ 请通过回复方式绑定")
				return
			}
			if GetBindNum(m.Chat.ID) == config.BindMaxNum {
				bot.Send(m.Chat, "⚠ 已经达到最大可绑定数")
				return
			}
			bot.Send(m.Chat, "正在绑定中……")
			err := BindUser(m, UserClientId[m.Chat.ID], UserClientSecret[m.Chat.ID])
			if err != nil {
				bot.Send(m.Chat, err.Error())
			} else {
				bot.Send(m.Chat, "✨ 绑定成功!")
			}
			UserStatus[m.Chat.ID] = USNone
		}
	}
}
func bTask(m *tb.Message) {
	for _, a := range config.Admins {
		if a == m.Chat.ID {
			SignTask()
			return
		}
	}
	bot.Send(m.Chat, "⚠ 只有Bot管理员才有权限执行此操作")
}
func bLog(m *tb.Message) {
	flag := 0
	for _, a := range config.Admins {
		if a == m.Chat.ID {
			flag = 1
		}
	}
	if flag == 0 {
		bot.Send(m.Chat, "⚠ 只有Bot管理员才有权限执行此操作")
		return
	}
	file := config.LogBasePath + "latest.log"
	logfile := &tb.Document{
		File:     tb.FromDisk(file),
		FileName: "latest.log",
		MIME:     "text/plain",
	}
	bot.Send(m.Chat, logfile)
}
