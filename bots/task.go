package bots

import (
	"fmt"
	"github.com/iyear/E5SubBot/config"
	"github.com/iyear/E5SubBot/model"
	"github.com/iyear/E5SubBot/task"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"time"
)

var errorTimes map[int]int
var signErr map[int64]int
var unbindUsers []int64
var msgSender *Sender

func InitTask() {
	errorTimes = make(map[int]int)
	msgSender = NewSender()

	c := cron.New()
	c.AddFunc(config.Cron, SignTask)
	c.Start()
}
func SignTask() {
	msgSender.Init(config.MaxGoroutines)

	signErr = make(map[int64]int)
	unbindUsers = nil

	var clients []*model.Client
	if result := model.DB.Find(&clients); result.Error != nil {
		zap.S().Errorw("failed to get all clients",
			"error", result.Error)
		return
	}

	fmt.Printf("clients: %d goroutines:%d\n",
		len(clients),
		config.MaxGoroutines,
	)

	start := time.Now()

	errClients := task.Sign(clients)

	for _, errClient := range errClients {
		if errClient.Err != nil {
			opErrorSign(errClient)
			continue
		}
		// 请求一次成功清零errorTimes，避免接口的偶然错误积累导致账号被清退
		errorTimes[errClient.ID] = 0
		model.DB.Save(&errClient.Client)
	}

	timeSpending := time.Since(start).Seconds()
	summarySignTaskForUsers(errClients)
	summarySignTaskForAdmins(errClients, timeSpending)

	msgSender.Stop()
}

func summarySignTaskForAdmins(errClients []*model.ErrClient, timeSpending float64) {
	var Count = len(errClients)
	var ErrCount int
	var ErrUserStr string
	var UnbindUserStr string
	for err, count := range signErr {
		ErrCount += count
		ErrUserStr += fmt.Sprintf("[%d](tg://user?id=%d)\n", err, err)
	}
	for _, unbindUser := range unbindUsers {
		UnbindUserStr += fmt.Sprintf("[%d](tg://user?id=%d)\n", unbindUser, unbindUser)
	}
	for _, admin := range config.Admins {
		a := admin
		msgSender.SendMessageByID(a, fmt.Sprintf("任务反馈(管理员)\n完成时间: %s\n用时: %.2fs\n结果: %d/%d\n错误账户: \n%s\n清退账户: \n%s",
			time.Now().Format("2006-01-02 15:04:05"),
			timeSpending,
			Count-ErrCount, Count,
			ErrUserStr, UnbindUserStr,
		),
			tb.ModeMarkdown,
		)
	}
}
func summarySignTaskForUsers(errClients []*model.ErrClient) {

	var isSent map[int64]bool
	isSent = make(map[int64]bool)

	for _, errClient := range errClients {
		errClient := errClient
		// pending SignErrNum
		if errorTimes[errClient.ID] > config.MaxErrTimes {
			if result := model.DB.Delete(&errClient.Client); result.Error != nil {
				zap.S().Errorw("failed to delete data",
					"error", result.Error,
					"id", errClient.ID,
				)
				continue
			}

			unbindUsers = append(unbindUsers, errClient.TgId)

			msgSender.SendMessageByID(errClient.TgId, fmt.Sprintf("您的账户因达到错误上限而被自动解绑\n后会有期!\n\n别名: %s\nclient_id: %s\nclient_secret: %s",
				errClient.Alias,
				errClient.ClientId,
				errClient.ClientSecret,
			))
			continue

		}
		if isSent[errClient.TgId] {
			continue
		}
		signOK := GetBindNum(errClient.TgId) - signErr[errClient.TgId]

		msgSender.SendMessageByID(errClient.TgId,
			fmt.Sprintf("任务反馈\n时间: %s\n结果:%d/%d",
				time.Now().Format("2006-01-02 15:04:05"),
				signOK,
				signErr[errClient.TgId]+signOK,
			),
		)
		isSent[errClient.TgId] = true
		time.Sleep(time.Millisecond * 100)
	}
}
func opErrorSign(errClient *model.ErrClient) {
	errorTimes[errClient.ID]++
	signErr[errClient.TgId]++

	UnBindBtn := tb.InlineButton{Unique: "un" + errClient.MsId, Text: "点击解绑", Data: strconv.Itoa(errClient.ID)}
	bot.Handle(&UnBindBtn, bUnBindInlineBtn)

	msgSender.SendMessageByID(errClient.TgId,
		fmt.Sprintf("您的帐户 %s 在执行时出现了错误\n您可以选择解绑该用户\n错误: %s",
			errClient.Alias, errClient.Err),
		&tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{UnBindBtn}}},
	)
}
