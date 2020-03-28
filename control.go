package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"time"
)

//If Successfully return "",else return error information
func BindUser(m *tb.Message) string {
	fmt.Printf("%d Begin Bind\n", m.Chat.ID)
	tmp := strings.Split(m.Text, " ")
	fmt.Println("alias: " + tmp[1])
	if len(tmp) != 2 {
		fmt.Printf("%d Bind error:Wrong Bind Format\n", m.Chat.ID)
		return "授权格式错误"
	}
	alias := tmp[1]
	code := GetURLValue(tmp[0], "code")
	fmt.Println(code)
	access, refresh := MSFirGetToken(code)
	if refresh == "" {
		fmt.Printf("%d Bind error:GetRefreshToken\n", m.Chat.ID)
		return "获取RefreshToken失败"
	}

	//token has gotten
	bot.Send(m.Chat, "Token获取成功!")
	info := MSGetUserInfo(access)
	fmt.Printf("TGID:%d Refresh Token: %s\n", m.Chat.ID, refresh)
	if info == "" {
		fmt.Printf("%d Bind error:Getinfo\n", m.Chat.ID)
		return "获取用户信息错误"
	}

	var u MSData
	u.tgId = m.Chat.ID
	u.refreshToken = refresh
	u.msId = gjson.Get(info, "id").String()
	u.uptime = time.Now()
	u.other = SetJsonValue("{}", "alias", alias)
	//MS User Is Exist
	if MSUserIsExist(u.tgId, u.msId) {
		fmt.Printf("%d Bind error:MSUserHasExisted\n", m.Chat.ID)
		return "该ID对应的用户已经绑定过了"
	}
	//MS information has gotten
	bot.Send(m.Chat, "MS_ID： "+u.msId+"\nuserPrincipalName： "+gjson.Get(info, "userPrincipalName").String()+"\ndisplayName： "+gjson.Get(info, "displayName").String()+"\n")
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
func MSUserIsExist(tgId int64, msId string) bool {
	data := QueryDataByTG(db, tgId)
	var res MSData
	for _, res = range data {
		if res.msId == msId {
			return true
		}
	}
	return false
}

//SignTask
func SignTask() {
	data := QueryDataAll(db)
	for _, u := range data {
		access := MSGetToken(u.refreshToken)
		if access == "" {
			fmt.Println(u.msId + "Sign ERROR:AccessTokenGet")
			continue
		}
		if !OutLookGetMails(access) {
			fmt.Println(u.msId + "Sign ERROR:ReadMails")
			continue
		}
		fmt.Println(u.msId + " Sign OK!")
		u.uptime = time.Now()
		if ok, err := UpdateData(db, u); !ok {
			fmt.Printf("%s Update Data ERROR: %s\n", u.msId, err)
		}
	}
}
