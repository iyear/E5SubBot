package util

import (
	"github.com/iyear/E5SubBot/app/bot/internal/config"
	tele "gopkg.in/telebot.v3"
)

func SendBatch(bot *tele.Bot, to []int64, what interface{}, opts ...interface{}) ([]*tele.Message, error) {
	msgs := make([]*tele.Message, 0)
	for _, id := range to {
		msg, err := bot.Send(&tele.Chat{ID: id}, what, opts...)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func EditBatch(bot *tele.Bot, msgs []*tele.Message, what interface{}, opts ...interface{}) error {
	for _, msg := range msgs {
		if _, err := bot.Edit(msg, what, opts...); err != nil {
			return err
		}
	}
	return nil
}

func IsAdmin(tid int64) bool {
	for _, a := range config.C.Bot.Admin {
		if a == tid {
			return true
		}
	}
	return false
}

func appendBack(back tele.InlineButton, opts ...interface{}) []interface{} {
	if len(opts) == 0 {
		return []interface{}{&tele.ReplyMarkup{InlineKeyboard: [][]tele.InlineButton{{back}}}}
	}

	for i, opt := range opts {
		switch t := opt.(type) {
		case *tele.SendOptions:
			if t.ReplyMarkup == nil {
				t.ReplyMarkup = &tele.ReplyMarkup{InlineKeyboard: [][]tele.InlineButton{{back}}}
			} else if t.ReplyMarkup.InlineKeyboard == nil {
				t.ReplyMarkup.InlineKeyboard = [][]tele.InlineButton{{back}}
			} else {
				t.ReplyMarkup.InlineKeyboard = append(t.ReplyMarkup.InlineKeyboard, []tele.InlineButton{back})
			}
			opts[i] = t
		case *tele.ReplyMarkup:
			if t.InlineKeyboard == nil {
				t.InlineKeyboard = [][]tele.InlineButton{{back}}
			} else {
				t.InlineKeyboard = append(t.InlineKeyboard, []tele.InlineButton{back})
			}
			opts[i] = t
		}
	}
	return opts
}

func doWithBack(c tele.Context, fn func(what interface{}, opts ...interface{}) error, what interface{}, opts ...interface{}) error {
	return fn(what, appendBack(GetScope(c).TMPL.B.Back, opts...)...)
}

func EditOrSendWithBack(c tele.Context, what interface{}, opts ...interface{}) error {
	return doWithBack(c, c.EditOrSend, what, opts...)
}

func SendWithBack(c tele.Context, what interface{}, opts ...interface{}) error {
	return doWithBack(c, c.Send, what, opts...)
}
