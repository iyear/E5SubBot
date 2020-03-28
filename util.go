package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func CheckErr(err error) bool {
	if err != nil {
		log.Println(err)
		fmt.Println("error: ", err.Error())
		panic(err)
		return false
	}
	return true
}
func FileExist(Path string) bool {
	if _, err := os.Stat(Path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			CheckErr(err)
		}
	}
	return true
}
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start)
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
func GetURLValue(Url, key string) string {
	u, _ := url.Parse(Url)
	query := u.Query()
	//fmt.Println(query.Get(key))
	return query.Get(key)
}

//return result
func SetJsonValue(sjson, key, value string) (RJson string) {
	var m map[string]interface{}
	m = make(map[string]interface{})
	err := json.Unmarshal([]byte(sjson), &m)
	CheckErr(err)
	m[key] = value
	data, _ := json.Marshal(m)
	return string(data)
}

//Returns a json
func MarshalMSData(u MSData) string {
	type MSDataOnlyString struct {
		TgId         string `json:"tgId"`
		RefreshToken string `json:"refreshToken"`
		MsId         string `json:"msId"`
		Uptime       string `json:"uptime"`
		Other        string `json:"other"`
	}
	var MSNoTime MSDataOnlyString
	MSNoTime.TgId = strconv.FormatInt(u.tgId, 10)
	MSNoTime.RefreshToken = u.refreshToken
	MSNoTime.MsId = u.msId
	MSNoTime.Uptime = u.uptime.Format("2006-01-02 15:04:05")
	MSNoTime.Other = u.other
	result, _ := json.Marshal(MSNoTime)
	return string(result)
}

//返回一个32位md5加密后的字符串
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

//返回一个16位md5加密后的字符串
func Get16MD5Encode(data string) string {
	return GetMD5Encode(data)[8:24]
}
