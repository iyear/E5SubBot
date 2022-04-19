package util

import (
	"github.com/iyear/E5SubBot/app/bot/internal/model"
	"github.com/iyear/E5SubBot/pkg/conf"
	tele "gopkg.in/telebot.v3"
)

func GetScope(c tele.Context) *model.Scope {
	return c.Get(conf.ContextScope).(*model.Scope)
}
