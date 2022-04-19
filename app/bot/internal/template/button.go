package template

import (
	"github.com/iyear/E5SubBot/app/bot/internal/config"
	"github.com/iyear/E5SubBot/app/bot/internal/model"
	tele "gopkg.in/telebot.v3"
)

var btns = make(map[string]*model.Button)

func bFrom(code string) *model.Button {
	if t, ok := btns[code]; ok {
		return t
	}
	return btns[config.C.Ctrl.DefaultLang]
}

func setButtons(code string, t *model.Template) {
	btns[code] = newButton(t.Button)
}

func newButton(t *model.ButtonTmpl) *model.Button {
	b := &model.Button{
		Back: tele.InlineButton{Unique: "back", Text: t.Back},
	}
	return b
}
