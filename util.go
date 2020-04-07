package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

//true=>no error
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func CheckErr(err error) bool {
	if err != nil {
		log.Println(err)
		fmt.Println("ERROR")
		panic(err)
		return false
	}
	return true
}
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
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

//只返回文件名
func GetPathFiles(path string) []string {
	files, _ := ioutil.ReadDir(path)
	var t []string
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			t = append(t, file.Name())
		}
	}
	return t
}

//输入文件夹路径，返回最近n个log的路径，不到n个返回所有
func GetRecentLogs(path string, n int) []string {
	var paths []string
	if !PathExists(path) {
		return paths
	}
	//path末尾检查/
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	data := time.Now()
	d, _ := time.ParseDuration("-24h")
	//不到n个返回所有
	nt := Min(n, len(GetPathFiles(path)))
	//fmt.Println(nt)
	for i := 1; i <= nt; {
		if PathExists(path + data.Format("2006-01-02") + ".log") {
			paths = append(paths, data.Format("2006-01-02")+".log")
			i++
		}
		data = data.Add(d)
	}
	return paths
}
