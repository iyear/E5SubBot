package bots

import (
	"encoding/json"
	"fmt"
	"github.com/iyear/E5SubBot/config"
	"github.com/iyear/E5SubBot/service/srv_client"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"os"
	"strconv"
)

func bStart(m *tb.Message) {
	bot.Send(m.Sender, config.WelcomeContent)
	bHelp(m)
}

func bExport(m *tb.Message) {
	type ClientExport struct {
		Alias        string `json:"alias"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RefreshToken string `json:"refresh_token"`
		Other        string `json:"other"`
	}
	var exports []*ClientExport
	data := srv_client.GetClients(m.Chat.ID)
	if len(data) == 0 {
		bot.Send(m.Chat, "⚠ 你还没有绑定过账户嗷~")
		return
	}
	for _, u := range data {
		var cExport = &ClientExport{
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
		fmt.Sprintf("%s\n%s", config.HelpContent, config.Notice),
		&tb.SendOptions{DisableWebPagePreview: false})
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
