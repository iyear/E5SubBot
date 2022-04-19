package handler

import (
	"fmt"
	"github.com/iyear/E5SubBot/app/bot/internal/middleware"
	"github.com/iyear/E5SubBot/app/bot/internal/template"
	"github.com/iyear/E5SubBot/app/bot/internal/util"
	"github.com/iyear/E5SubBot/pkg/conf"
	"github.com/iyear/E5SubBot/pkg/format"
	"github.com/iyear/E5SubBot/pkg/utils"
	iso6391 "github.com/iyear/iso-639-1"
	"github.com/patrickmn/go-cache"
	"go.etcd.io/bbolt"
	tele "gopkg.in/telebot.v3"
)

func SettingsStart(c tele.Context) error {
	sp := util.GetScope(c)
	langBtn := sp.TMPL.B.SettingsLanguage

	return util.EditOrSendWithBack(c, sp.TMPL.I.Settings.Desc.T(nil),
		&tele.ReplyMarkup{InlineKeyboard: [][]tele.InlineButton{{langBtn}}})
}

func SettingsLanguage(c tele.Context) error {
	sp := util.GetScope(c)

	langBtns := make([][]tele.InlineButton, 0)

	nowLang := util.GetLangCode(sp.DB, c.Chat().ID)

	for _, code := range template.Langs() {
		langBtn := sp.TMPL.B.SettingsLanguagePlain
		lang := iso6391.FromCode(code)

		langBtn.Text = fmt.Sprintf("%s%s (%s)", utils.IF(nowLang == code, "✅ ", ""), lang.Name, lang.NativeName)
		langBtn.Data = code
		langBtns = append(langBtns, []tele.InlineButton{langBtn})
	}

	return util.EditOrSendWithBack(c, sp.TMPL.I.Settings.Language.T(iso6391.FromCode(nowLang).NativeName),
		&tele.ReplyMarkup{InlineKeyboard: langBtns})
}

func SettingsSetLanguage(c tele.Context) error {
	sp := util.GetScope(c)

	// 相同则不做任何事
	if util.GetLangCode(sp.DB, c.Chat().ID) == c.Data() {
		return nil
	}

	err := sp.DB.KV.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket([]byte(conf.BucketLanguage)).Put(format.Key.BoltLanguage(c.Chat().ID), []byte(c.Data())); err != nil {
			return err
		}

		// 刷新缓存
		sp.DB.Cache.Set(format.Key.CacheLanguage(c.Chat().ID), c.Data(), cache.NoExpiration)
		return nil
	})
	if err != nil {
		return err
	}

	// 手动走一遍中间件刷新当前页面的语言
	return middleware.SetScope(sp)(SettingsLanguage)(c)
}
