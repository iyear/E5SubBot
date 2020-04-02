package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
	"time"
)

var SignOk map[int64]int

//If Successfully return "",else return error information
func BindUser(m *tb.Message, cid, cse string) string {
	fmt.Printf("%d Begin Bind\n", m.Chat.ID)
	tmp := strings.Split(m.Text, " ")
	if len(tmp) != 2 {
		fmt.Printf("%d Bind error:Wrong Bind Format\n", m.Chat.ID)
		return "授权格式错误"
	}
	fmt.Println("alias: " + tmp[1])
	alias := tmp[1]
	code := GetURLValue(tmp[0], "code")
	//fmt.Println(code)
	access, refresh := MSFirGetToken(code, cid, cse)
	if refresh == "" {
		fmt.Printf("%d Bind error:GetRefreshToken\n", m.Chat.ID)
		return "获取RefreshToken失败"
	}

	//token has gotten
	bot.Send(m.Chat, "Token获取成功!")
	info := MSGetUserInfo(access)
	//fmt.Printf("TGID:%d Refresh Token: %s\n", m.Chat.ID, refresh)
	if info == "" {
		fmt.Printf("%d Bind error:Getinfo\n", m.Chat.ID)
		return "获取用户信息错误"
	}

	var u MSData
	u.tgId = m.Chat.ID
	u.refreshToken = refresh
	//TG的Data传递最高64bytes,一些msid超过了报错BUTTON_DATA_INVALID (0)，采取md5
	u.msId = Get16MD5Encode(gjson.Get(info, "id").String())
	u.uptime = time.Now().Unix()
	fmt.Println(u.uptime)
	u.alias = alias
	u.clientId = cid
	u.clientSecret = cse
	u.other = ""
	//MS User Is Exist
	if MSAppIsExist(u.tgId, u.clientId) {
		fmt.Printf("%d Bind error:MSUserHasExisted\n", m.Chat.ID)
		return "该应用已经绑定过了，无需重复绑定"
	}
	//MS information has gotten
	bot.Send(m.Chat, "MS_ID(MD5)： "+u.msId+"\nuserPrincipalName： "+gjson.Get(info, "userPrincipalName").String()+"\ndisplayName： "+gjson.Get(info, "displayName").String()+"\n")
	if ok, err := AddData(db, u); !ok {
		fmt.Printf("%d Bind error: %s\n", m.Chat.ID, err)
		return "数据库写入错误"
	}
	fmt.Printf("%d Bind Successfully!\n", m.Chat.ID)
	return ""
}

//get bind num
func GetBindNum(tgId int64) int {
	data := QueryDataByTG(db, tgId)
	return len(data)
}

//return true => exist
func MSAppIsExist(tgId int64, clientId string) bool {
	data := QueryDataByTG(db, tgId)
	var res MSData
	for _, res = range data {
		if res.clientId == clientId {
			return true
		}
	}
	return false
}

//SignTask
func SignTask() {
	var SignOk map[int64]int
	var SignErr []string
	var num, signOk int
	SignOk = make(map[int64]int)
	fmt.Println("----Task Begin----")
	fmt.Println("Time:" + time.Now().Format("2006-01-02 15:04:05"))
	data := QueryDataAll(db)
	num = len(data)
	fmt.Println("Start Sign")
	for _, u := range data {
		e := ""
		pre := "您的账户:" + u.alias + "\n在任务执行时出现了错误!\n错误:"
		access := MSGetToken(u.refreshToken, u.clientId, u.clientSecret)
		chat, _ := bot.ChatByID(strconv.FormatInt(u.tgId, 10))
		if access == "" {
			e = "Sign ERROR:GetAccessToken"
			fmt.Println(u.msId + e)
			bot.Send(chat, pre+e)
			SignErr = append(SignErr, u.msId)
			continue
		}
		if !OutLookGetMails(access) {
			e = "Sign ERROR:ReadMails"
			fmt.Println(u.msId + " Sign ERROR:ReadMails")
			bot.Send(chat, pre+e)
			SignErr = append(SignErr, u.msId)
			continue
		}
		u.uptime = time.Now().Unix()
		if ok, err := UpdateData(db, u); !ok {
			e = "Update Data ERROR:"
			fmt.Printf("%s Update Data ERROR: %s\n", u.msId, err)
			bot.Send(chat, pre+e)
			SignErr = append(SignErr, u.msId)
			continue
		}
		fmt.Println(u.msId + " Sign OK!")
		SignOk[u.tgId]++
		signOk++
	}
	fmt.Println("Sign End,Start Send")
	var isSend map[int64]bool
	isSend = make(map[int64]bool)
	//用户任务反馈
	for _, u := range data {
		if !isSend[u.tgId] {
			chat, err := bot.ChatByID(strconv.FormatInt(u.tgId, 10))
			if err != nil {
				fmt.Println("Send Result ERROR")
				continue
			}
			bot.Send(chat, "任务反馈\n时间: "+time.Now().Format("2006-01-02 15:04:05")+"\n结果: "+strconv.Itoa(SignOk[u.tgId])+"/"+strconv.Itoa(GetBindNum(u.tgId)))
			isSend[u.tgId] = true
		}
	}
	//管理员任务反馈
	var ErrUser string
	for _, eu := range SignErr {
		ErrUser = ErrUser + eu + "\n"
	}
	for _, a := range admin {
		chat, _ := bot.ChatByID(strconv.FormatInt(a, 10))
		bot.Send(chat, "任务反馈(管理员)\n完成时间: "+time.Now().Format("2006-01-02 15:04:05")+"\n结果: "+strconv.Itoa(signOk)+"/"+strconv.Itoa(num)+"\n错误账户msid:\n"+ErrUser)
	}
	fmt.Println("----Task End----")
}
func GetAdmin() []int64 {
	var result []int64
	admins := strings.Split(viper.GetString("admin"), ",")
	for _, v := range admins {
		id, _ := strconv.ParseInt(v, 10, 64)
		result = append(result, id)
	}
	fmt.Println(result)
	return result
}
