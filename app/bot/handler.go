package bot

import (
	"github.com/iyear/E5SubBot/app/bot/internal/handler"
	"github.com/iyear/E5SubBot/app/bot/internal/middleware"
	"github.com/iyear/E5SubBot/app/bot/internal/template"
	tele "gopkg.in/telebot.v3"
)

func makeHandlers(bot *tele.Bot) {
	b := template.From("").B
	h := bot.Group()
	h.Use(middleware.Private())
	{
		h.Handle("/start", handler.OnStart)

		h.Handle(&b.Back, handler.OnStart)
		h.Handle(&b.StartMy, handler.My)
	}

}
