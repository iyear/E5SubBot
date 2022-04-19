package model

import (
	tele "gopkg.in/telebot.v3"
	"strings"
	"text/template"
)

type Template struct {
	Button *ButtonTmpl `mapstructure:"button"`
	Start  struct {
		Desc Tmpl `mapstructure:"desc"`
	} `mapstructure:"start"`
}

type ButtonTmpl struct {
	Back        string `mapstructure:"back"`
	StartBind   string `mapstructure:"start_bind"`
	StartUnbind string `mapstructure:"start_unbind"`
	StartExport string `mapstructure:"start_export"`
	StartImport string `mapstructure:"start_import"`
	StartMy     string `mapstructure:"start_my"`
	StartHelp   string `mapstructure:"start_help"`
}

type Button struct {
	Back tele.InlineButton

	StartBind   tele.InlineButton
	StartUnbind tele.InlineButton
	StartExport tele.InlineButton
	StartImport tele.InlineButton
	StartMy     tele.InlineButton
	StartHelp   tele.InlineButton
}

type Tmpl struct {
	tmpl *template.Template
}

func NewTmpl(t *template.Template) Tmpl {
	return Tmpl{tmpl: t}
}

func (t Tmpl) T(data interface{}) string {
	if t.tmpl == nil {
		return ""
	}
	var buf strings.Builder
	err := t.tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
