package handler

import (
	"github.com/iyear/E5SubBot/app/bot/internal/config"
	"github.com/iyear/E5SubBot/app/bot/internal/template"
	"github.com/iyear/E5SubBot/app/bot/internal/util"
	tele "gopkg.in/telebot.v3"
)

func OnStart(c tele.Context) error {
	sp := util.GetScope(c)
	if c.Message().Payload == "" {
		chat := c.Chat()
		b := sp.TMPL.B
		return c.EditOrSend(sp.TMPL.I.Start.Desc.T(&template.MStartWelcome{
			ID:       chat.ID,
			Username: chat.Username,
			Notice:   config.C.Biz.Notice,
		}), &tele.SendOptions{
			DisableWebPagePreview: true,
			ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: [][]tele.InlineButton{
				{b.StartMy, b.StartHelp},
				{b.StartBind, b.StartUnbind},
				{b.StartExport, b.StartImport},
			}},
		})
	}

	return nil
}
