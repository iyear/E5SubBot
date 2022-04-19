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
		Back:          tele.InlineButton{Unique: "back", Text: t.Back},
		StartBind:     tele.InlineButton{Unique: "start_bind", Text: t.StartBind},
		StartUnbind:   tele.InlineButton{Unique: "start_unbind", Text: t.StartUnbind},
		StartExport:   tele.InlineButton{Unique: "start_export", Text: t.StartExport},
		StartImport:   tele.InlineButton{Unique: "start_import", Text: t.StartImport},
		StartMy:       tele.InlineButton{Unique: "start_my", Text: t.StartMy},
		StartSettings: tele.InlineButton{Unique: "start_settings", Text: t.StartSettings},

		SettingsLanguage:      tele.InlineButton{Unique: "settings_language", Text: t.SettingsLanguage},
		SettingsLanguagePlain: tele.InlineButton{Unique: "settings_set_language"},
	}
	return b
}
