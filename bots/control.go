package bots

import (
	"errors"
	"fmt"
	"github.com/iyear/E5SubBot/model"
	"github.com/iyear/E5SubBot/util"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
	"time"
)

var SignOk map[int64]int

// BindUser If Successfully return "",else return error information
func BindUser(m *tb.Message, ClientId, ClientSecret string) error {
	tmp := strings.Split(m.Text, " ")
	if len(tmp) != 2 {
		return errors.New("wrong format")
	}
	code := util.GetURLValue(tmp[0], "code")
	Alias := tmp[1]
	cli := model.NewClient(ClientId, ClientSecret)
	if err := cli.GetTokenWithCode(code); err != nil {
		return err
	}
	bot.Send(m.Chat, "Token获取成功!")

	info, err := cli.GetUserInfo()
	if err != nil {
		return err
	}
	var u = &model.Client{
		TgId: m.Chat.ID,
		//TG的Data传递最高64bytes,一些MsId超过了报错BUTTON_DATA_INVALID (0)，采取md5
		RefreshToken: cli.RefreshToken,
		MsId:         util.Get16MD5Encode(gjson.Get(info, "id").String()),
		Alias:        Alias,
		ClientId:     ClientId,
		ClientSecret: ClientSecret,
		Other:        "",
	}

	//MS User Is Exist
	if MSAppIsExist(u.TgId, u.ClientId) {
		return errors.New("该应用已经绑定过了，无需重复绑定")
	}
	//MS information has gotten
	bot.Send(m.Chat,
		fmt.Sprintf("MS_ID(MD5)： %s\nuserPrincipalName： %s\ndisplayName： %s\n",
			u.MsId,
			gjson.Get(info, "userPrincipalName").String(),
			gjson.Get(info, "displayName").String()),
	)

	if result := model.DB.Create(&u); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetBindNum get bind num
func GetBindNum(TgId int64) int {
	var bindings []*model.Client
	result := model.DB.Where("tg_id = ?", TgId).Find(&bindings)
	return int(result.RowsAffected)
}

// MSAppIsExist return true => exist
func MSAppIsExist(TgId int64, ClientId string) bool {
	result := model.DB.
		Where("tg_id = ? AND client_id = ?", TgId, ClientId).
		First(&model.Client{})
	return util.IF(result.RowsAffected == 0, false, true).(bool)
}

// func SignTask() {
//	var (
//		SignOk      map[int64]int
//		SignErr     []string
//		UnbindUser  []string
//		num, signOk int
//	)
//	SignOk = make(map[int64]int)
//	fmt.Println("----Task Begin----")
//	fmt.Println("Time:" + time.Now().Format("2006-01-02 15:04:05"))
//	data := model.QueryDataAll()
//	num = len(data)
//	fmt.Println("Start Sign")
//	//签到任务
//	for _, u := range data {
//		pre := "您的账户: " + u.Alias + "\n在任务执行时出现了错误!\n错误:"
//		chat, err := bot.ChatByID(strconv.FormatInt(u.TgId, 10))
//		if err != nil {
//			zap.S().Errorw("wrong chat id", "error", err, "tg_id", u.TgId)
//			continue
//		}
//		//生成解绑按钮
//		var inlineKeys [][]tb.InlineButton
//		UnBindBtn := tb.InlineButton{Unique: "un" + u.MsId, Text: "点击解绑该账户", Data: u.MsId}
//		bot.Handle(&UnBindBtn, bUnBindInlineBtn)
//		inlineKeys = append(inlineKeys, []tb.InlineButton{UnBindBtn})
//		tmpBtn := &tb.ReplyMarkup{InlineKeyboard: inlineKeys}
//
//		se := u.MsId + " ( @" + chat.Username + " )"
//		client := model.NewClient(u.ClientId, u.ClientSecret)
//		if err := client.GetOutlookMails(); err != nil {
//			zap.S().Errorw("failed to get outlook mails", "error", err, "ms_id", u.MsId)
//			bot.Send(chat, pre+err.Error(), tmpBtn)
//			ErrorTimes[u.MsId]++
//			SignErr = append(SignErr, se)
//			continue
//		}
//		u.Uptime = time.Now().Unix()
//		u.RefreshToken = newRefreshToken
//		if ok, err := model.UpdateData(u); !ok {
//			zap.S().Errorw("failed to update db data", "error", err, "ms_id", u.MsId)
//			bot.Send(chat, pre+err.Error(), tmpBtn)
//			SignErr = append(SignErr, se)
//			ErrorTimes[u.MsId]++
//			continue
//		}
//		fmt.Println(u.MsId + " Sign OK!")
//		SignOk[u.TgId]++
//		signOk++
//	}
//	fmt.Println("Sign End,Start Send")
//	var isSend map[int64]bool
//	isSend = make(map[int64]bool)
//	//用户任务反馈
//	for _, u := range data {
//		chat, err := bot.ChatByID(strconv.FormatInt(u.TgId, 10))
//		if err != nil {
//			zap.S().Errorw("failed to get chat", "error", err, "tg_id", u.TgId)
//			continue
//		}
//		//错误上限账户清退
//		if ErrorTimes[u.MsId] == ErrMaxTimes {
//			zap.S().Errorw("binding max num limit", "ms_id", u.MsId)
//			if ok, err := model.DelData(u.MsId); !ok {
//				zap.S().Errorw("failed to delete db data", "error", err, "ms_id", u.MsId)
//			} else {
//				UnbindUser = append(UnbindUser, u.MsId+" ( @"+chat.Username+" )")
//				bot.Send(chat, "您的账户因达到错误上限而被自动解绑\n后会有期!\n\n别名: "+u.Alias+"\nclient_id: "+u.ClientId+"\nclient_secret: "+u.ClientSecret)
//			}
//
//		}
//		if !isSend[u.TgId] {
//			//静默发送，过多消息很烦
//			bot.Send(chat, "任务反馈\n时间: "+time.Now().Format("2006-01-02 15:04:05")+"\n结果: "+strconv.Itoa(SignOk[u.TgId])+"/"+strconv.Itoa(GetBindNum(u.TgId)), &tb.SendOptions{DisableNotification: true})
//			isSend[u.TgId] = true
//		}
//	}
//	//管理员任务反馈
//	var ErrUserStr string
//	var UnbindUserStr string
//	for _, eu := range SignErr {
//		ErrUserStr = ErrUserStr + eu + "\n"
//	}
//	for _, ubu := range UnbindUser {
//		UnbindUserStr = UnbindUserStr + ubu + "\n"
//	}
//	for _, a := range admin {
//		chat, err := bot.ChatByID(strconv.FormatInt(a, 10))
//		if err != nil {
//			zap.S().Errorw("failed to get chat", "error", err, "tg_id", a)
//			continue
//		}
//		bot.Send(chat, "任务反馈(管理员)\n完成时间: "+time.Now().Format("2006-01-02 15:04:05")+"\n结果: "+strconv.Itoa(signOk)+"/"+strconv.Itoa(num)+"\n错误账户:\n"+ErrUserStr+"\n清退账户:\n"+UnbindUserStr)
//	}
//	fmt.Println("----Task End----")
//}

func GetAdmins() []int64 {
	var result []int64
	admins := strings.Split(viper.GetString("admin"), ",")
	for _, v := range admins {
		id, _ := strconv.ParseInt(v, 10, 64)
		result = append(result, id)
	}
	return result
}
