package util

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/url"
	"os"
	"time"
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

func GetRecentLogs(path string, n int) []string {
	var paths []string
	if !PathExists(path) {
		return paths
	}

	if path[len(path)-1:] != "/" {
		path += "/"
	}
	data := time.Now()
	d, _ := time.ParseDuration("-24h")
	//不到n个返回所有
	nt := Min(n, len(GetPathFiles(path)))
	for i := 1; i <= nt; {
		if PathExists(path + data.Format("2006-01-02") + ".log") {
			paths = append(paths, data.Format("2006-01-02")+".log")
			i++
		}
		data = data.Add(d)
	}
	return paths
}
func IF(f bool, a interface{}, b interface{}) interface{} {
	if f {
		return a
	}
	return b
}
