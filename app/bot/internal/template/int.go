package template

import (
	_ "embed"
	"github.com/iyear/E5SubBot/app/bot/internal/config"
	"github.com/iyear/E5SubBot/app/bot/internal/model"
)

var (
	ints  = make(map[string]*model.Template)
	codes []string
)

func iFrom(code string) *model.Template {
	if t, ok := ints[code]; ok {
		return t
	}
	return ints[config.C.Ctrl.DefaultLang]
}

// Langs 语言种类以internal为准
func Langs() []string {
	return codes
}
