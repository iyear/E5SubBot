package model

import (
	"github.com/patrickmn/go-cache"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Scope struct {
	DB   *DB
	TMPL *TMPL
	Log  *zap.SugaredLogger
}

type TMPL struct {
	I *Template
	B *Button
}

type DB struct {
	DB    *gorm.DB
	KV    *bbolt.DB
	Cache *cache.Cache
}
