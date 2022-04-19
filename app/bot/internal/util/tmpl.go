package util

import (
	"github.com/iyear/E5SubBot/app/bot/internal/config"
	"github.com/iyear/E5SubBot/app/bot/internal/model"
	"github.com/iyear/E5SubBot/pkg/conf"
	"github.com/iyear/E5SubBot/pkg/format"
	"go.etcd.io/bbolt"
	"time"
)

func GetLangCode(db *model.DB, tid int64) string {
	t, found := db.Cache.Get(format.Key.CacheLanguage(tid))
	if found {
		return t.(string)
	}
	code := ""

	_ = db.KV.View(func(tx *bbolt.Tx) error {
		if v := tx.Bucket([]byte(conf.BucketLanguage)).Get(format.Key.BoltLanguage(tid)); v != nil {
			code = string(v)
			return nil
		}
		code = config.C.Ctrl.DefaultLang
		return nil
	})

	// 缓存
	db.Cache.Set(format.Key.CacheLanguage(tid), code, 15*time.Minute)
	return code
}
