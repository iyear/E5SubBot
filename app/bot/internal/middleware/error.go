package middleware

import (
	"github.com/iyear/E5SubBot/app/bot/internal/util"
	tele "gopkg.in/telebot.v3"
)

func OnError() func(err error, ctx tele.Context) {
	return func(err error, ctx tele.Context) {
		if err != nil {
			sp := util.GetScope(ctx)
			if len(ctx.Recipient().Recipient()) > 0 {
				sp.Log.Errorw("error",
					"err", err,
					"recipient", ctx.Recipient().Recipient())
				return
			}
			sp.Log.Errorw("error", "err", err)
		}
	}
}
