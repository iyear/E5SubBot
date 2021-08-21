package util

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/url"
	"os"
)

// Min true=>no error
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
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
func GetURLValue(Url, key string) string {
	u, _ := url.Parse(Url)
	query := u.Query()
	//fmt.Println(query.Get(key))
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

// GetPathFiles only return file name
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
func IF(f bool, a interface{}, b interface{}) interface{} {
	if f {
		return a
	}
	return b
}
