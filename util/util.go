package util

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
)

func GetURLValue(Url, key string) string {
	u, _ := url.Parse(Url)
	query := u.Query()
	return query.Get(key)
}

func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Get16MD5Encode(data string) string {
	return GetMD5Encode(data)[8:24]
}
