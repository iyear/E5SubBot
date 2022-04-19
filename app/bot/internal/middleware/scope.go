package middleware

import (
	"github.com/google/uuid"
	"github.com/iyear/E5SubBot/app/bot/internal/model"
	"github.com/iyear/E5SubBot/app/bot/internal/template"
	"github.com/iyear/E5SubBot/app/bot/internal/util"
	"github.com/iyear/E5SubBot/pkg/conf"
	tele "gopkg.in/telebot.v3"
)

func SetScope(sp *model.Scope) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			tid := getTID(c)
			lang := util.GetLangCode(sp.DB, tid)
			tmpl := template.From(lang)

			trace := uuid.New().String()
			c.Set(conf.ContextScope, &model.Scope{
				DB:   sp.DB,
				TMPL: tmpl,
				Log:  sp.Log.With("_trace", trace),
			})
			c.Set(conf.ContextTrace, trace)
			return next(c)
		}
	}
}

func getTID(c tele.Context) int64 {
	if c.Chat() != nil {
		return c.Chat().ID
	}
	if c.Query() != nil {
		return c.Query().Sender.ID
	}
	if c.Message() != nil {
		return c.Message().Sender.ID
	}
	return 0
}
