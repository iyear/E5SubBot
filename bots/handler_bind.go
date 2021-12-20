package bots

import (
	"fmt"
	"github.com/iyear/E5SubBot/config"
	"github.com/iyear/E5SubBot/model"
	"github.com/iyear/E5SubBot/pkg/microsoft"
	"github.com/iyear/E5SubBot/service/srv_client"
	"github.com/iyear/E5SubBot/util"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

func bBind(m *tb.Message) {
	bot.Send(m.Chat,
		fmt.Sprintf("👉 应用注册： [点击直达](%s)", microsoft.GetRegURL()),
		tb.ModeMarkdown,
	)

	bot.Send(m.Chat,
		"⚠ 请回复 `client_id(空格)client_secret`",
		&tb.SendOptions{ParseMode: tb.ModeMarkdown,
			ReplyMarkup: &tb.ReplyMarkup{ForceReply: true}},
	)

	UserStatus[m.Chat.ID] = StatusBind1
	UserClientId[m.Chat.ID] = m.Text
}

func bBind1(m *tb.Message) {
	if !m.IsReply() {
		bot.Send(m.Chat, "⚠ 请通过回复方式绑定")
		return
	}
	tmp := strings.Split(m.Text, " ")
	if len(tmp) != 2 {
		bot.Send(m.Chat, "⚠ 错误的格式")
		return
	}
	id := tmp[0]
	secret := tmp[1]
	bot.Send(m.Chat,
		fmt.Sprintf("👉 授权账户： [点击直达](%s)", microsoft.GetAuthURL(id)),
		tb.ModeMarkdown,
	)

	bot.Send(m.Chat,
		"⚠ 请回复`http://localhost/……(空格)别名`(用于管理)",
		&tb.SendOptions{ParseMode: tb.ModeMarkdown,
			ReplyMarkup: &tb.ReplyMarkup{ForceReply: true},
		},
	)
	UserStatus[m.Chat.ID] = StatusBind2
	UserClientId[m.Chat.ID] = id
	UserClientSecret[m.Chat.ID] = secret
}

func bBind2(m *tb.Message) {
	if !m.IsReply() {
		bot.Send(m.Chat, "⚠ 请通过回复方式绑定")
		return
	}
	if len(srv_client.GetClients(m.Chat.ID)) == config.BindMaxNum {
		bot.Send(m.Chat, "⚠ 已经达到最大可绑定数")
		return
	}
	bot.Send(m.Chat, "正在绑定中……")

	tmp := strings.Split(m.Text, " ")
	if len(tmp) != 2 {
		bot.Send(m.Chat, "😥 错误的格式")
	}
	code := util.GetURLValue(tmp[0], "code")
	alias := tmp[1]

	id := UserClientId[m.Chat.ID]
	secret := UserClientSecret[m.Chat.ID]

	refresh, err := microsoft.GetTokenWithCode(id, secret, code)
	if err != nil {
		bot.Send(m.Chat, fmt.Sprintf("无法获取RefreshToken ERROR:%s", err))
		return
	}
	bot.Send(m.Chat, "🎉 Token获取成功!")

	refresh, info, err := microsoft.GetUserInfo(id, secret, refresh)
	if err != nil {
		bot.Send(m.Chat, fmt.Sprintf("无法获取用户信息 ERROR:%s", err))
		return
	}
	c := &model.Client{
		TgId:         m.Chat.ID,
		RefreshToken: refresh,
		MsId:         util.Get16MD5Encode(gjson.Get(info, "id").String()),
		Alias:        alias,
		ClientId:     id,
		ClientSecret: secret,
		Other:        "",
	}

	if srv_client.IsExist(c.TgId, c.ClientId) {
		bot.Send(m.Chat, "⚠ 该应用已经绑定过了，无需重复绑定")
		return
	}

	bot.Send(m.Chat,
		fmt.Sprintf("ms_id：%s\nuserPrincipalName：%s\ndisplayName：%s",
			c.MsId,
			gjson.Get(info, "userPrincipalName").String(),
			gjson.Get(info, "displayName").String(),
		),
	)

	if err = srv_client.Add(c); err != nil {
		bot.Send(m.Chat, "😥 用户写入数据库失败")
		return
	}

	bot.Send(m.Chat, "✨ 绑定成功!")
	delete(UserStatus, m.Chat.ID)
	delete(UserClientId, m.Chat.ID)
	delete(UserClientSecret, m.Chat.ID)
}

func bUnBind(m *tb.Message) {
	var inlineKeys [][]tb.InlineButton
	clients := srv_client.GetClients(m.Chat.ID)

	for _, u := range clients {
		inlineBtn := tb.InlineButton{
			Unique: "unbind" + strconv.Itoa(u.ID),
			Text:   u.Alias,
			Data:   strconv.Itoa(u.ID),
		}
		bot.Handle(&inlineBtn, bUnBindInlineBtn)
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}

	bot.Send(m.Chat,
		fmt.Sprintf("⚠ 选择一个账户将其解绑\n\n当前绑定数: %d/%d", len(srv_client.GetClients(m.Chat.ID)), config.BindMaxNum),
		&tb.ReplyMarkup{InlineKeyboard: inlineKeys},
	)
}
func bUnBindInlineBtn(c *tb.Callback) {
	id, _ := strconv.Atoi(c.Data)
	if err := srv_client.Del(id); err != nil {
		zap.S().Errorw("failed to delete db data",
			"error", err,
			"id", c.Data,
		)
		bot.Send(c.Message.Chat, "⚠ 解绑失败!")
		return
	}
	bot.Send(c.Message.Chat, "✨ 解绑成功!")
	bot.Respond(c)
}
