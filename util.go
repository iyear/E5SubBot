package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	bStartContent string = "欢迎使用E5SubBot!"
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
