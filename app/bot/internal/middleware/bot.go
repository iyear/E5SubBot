package middleware

import (
	tele "gopkg.in/telebot.v3"
	"strings"
)

func AutoResponder() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Callback() != nil {
				defer func(c tele.Context) {
					_ = c.Respond()
				}(c)
			}
			return next(c) // continue execution chain
		}
	}
}

func Private() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Chat() != nil && c.Chat().Type != tele.ChatPrivate {
				return nil
			}
			return next(c)
		}
	}
}

func DeleteMsg() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			// query 发的消息不能删
			// 管理员命令提示不删
			if c.Callback() == nil && c.Message() != nil && !strings.HasPrefix(c.Message().Text, "/cmd") {
				if err := c.Delete(); err != nil {
					return err
				}
			}
			return next(c)
		}
	}
}
