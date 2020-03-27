package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func BindUser(m *tb.Message) bool {
	fmt.Printf("%d Begin Bind", m.Chat.ID)
	code := GetURLValue(m.Text, "code")
	fmt.Println(code)
	access, refresh := MSFirGetToken(code)
	if refresh == "" {
		fmt.Printf("%d Bind error:Getinfo", m.Chat.ID)
		return false
	}
	fmt.Printf("TGID:%d Refresh Token: %s", m.Chat.ID, refresh)
	//
	bot.Send(m.Chat, "Token获取成功!")
	info := MSGetUserInfo(access)
	fmt.Printf("TGID:%d Refresh Token: %s", m.Chat.ID, refresh)
	if info == "" {
		fmt.Printf("%d Bind error:Getinfo", m.Chat.ID)
		return false
	}
	var u MSData
	u.tgId = m.Chat.ID
	u.refreshToken = refresh
	u.msId = gjson.Get(info, "id").String()
	u.uptime = time.Now()
	u.other = ""
	//
	bot.Send(m.Chat, "MS_ID:"+u.msId+"\nuserPrincipalName: "+gjson.Get(info, "userPrincipalName").String()+"\ndisplayName"+gjson.Get(info, "displayName").String())
	if ok, err := AddData(db, u); !ok {
		fmt.Printf("%d Bind error: %s", m.Chat.ID, err)
		return false
	}
	fmt.Printf("%d Bind Successfully!", m.Chat.ID)
	return true
}
