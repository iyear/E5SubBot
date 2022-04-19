package model

import (
	tele "gopkg.in/telebot.v3"
	"strings"
	"text/template"
)

type Template struct {
	Button *ButtonTmpl `mapstructure:"button"`
}

type ButtonTmpl struct {
	Back string `mapstructure:"back"`
}

type Button struct {
	Back tele.InlineButton
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
