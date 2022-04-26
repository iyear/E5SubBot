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
	Settings struct {
		Desc     Tmpl `mapstructure:"desc"`
		Language Tmpl `mapstructure:"language"`
	} `mapstructure:"settings"`
	My struct {
		Desc Tmpl `mapstructure:"desc"`
		View Tmpl `mapstructure:"view"`
	}
}

type ButtonTmpl struct {
	Back          string `mapstructure:"back"`
	StartBind     string `mapstructure:"start_bind"`
	StartUnbind   string `mapstructure:"start_unbind"`
	StartExport   string `mapstructure:"start_export"`
	StartImport   string `mapstructure:"start_import"`
	StartMy       string `mapstructure:"start_my"`
	StartSettings string `mapstructure:"start_settings"`

	SettingsLanguage string `mapstructure:"settings_language"`
}

type Button struct {
	Back tele.InlineButton

	StartBind     tele.InlineButton
	StartUnbind   tele.InlineButton
	StartExport   tele.InlineButton
	StartImport   tele.InlineButton
	StartMy       tele.InlineButton
	StartSettings tele.InlineButton

	SettingsLanguage      tele.InlineButton
	SettingsLanguagePlain tele.InlineButton

	MyViewClient tele.InlineButton
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
