package handler

import (
	"github.com/iyear/E5SubBot/app/bot/internal/config"
	"github.com/iyear/E5SubBot/app/bot/internal/template"
	"github.com/iyear/E5SubBot/app/bot/internal/util"
	"github.com/iyear/E5SubBot/pkg/models"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"time"
)

func MyStart(c tele.Context) error {
	clients := make([]*models.Client, 0)
	sp := util.GetScope(c)

	if err := sp.DB.DB.Where("tg_id = ?", c.Chat().ID).Find(&clients).Error; err != nil {
		return err
	}

	var inlineKeys [][]tele.InlineButton
	for _, client := range clients {
		b := sp.TMPL.B.MyViewClient
		b.Text = client.Alias
		b.Data = strconv.FormatInt(client.ID, 10)
		inlineKeys = append(inlineKeys, []tele.InlineButton{b})
	}

	return util.EditOrSendWithBack(c,
		sp.TMPL.I.My.Desc.T(&template.MMyDesc{
			Current: len(clients),
			BindMax: config.C.Ctrl.BindMax,
		}),
		&tele.ReplyMarkup{InlineKeyboard: inlineKeys})
}

func MyViewClient(c tele.Context) error {
	sp := util.GetScope(c)

	id, err := strconv.Atoi(c.Data())
	if err != nil {
		return err
	}

	client := &models.Client{}
	if err = sp.DB.DB.Where("id = ?", id).First(&client).Error; err != nil {
		return err
	}

	return util.EditOrSendWithBack(c,
		sp.TMPL.I.My.View.T(&template.MMyView{
			Alias:        client.Alias,
			ClientID:     client.ClientId,
			ClientSecret: client.ClientSecret,
			UpdatedAt:    time.Unix(client.Uptime, 0).Format("2006-01-02 15:04:05"),
		}))
}
