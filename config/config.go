package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

const (
	LogBasePath    string = "./log/"
	WelcomeContent string = "欢迎使用E5SubBot!"
	HelpContent    string = `
	命令：
	/my 查看已绑定账户信息
	/bind  绑定新账户
	/unbind 解绑账户
	/export 导出账户信息(JSON)
	/help 帮助
	源码及使用方法：https://github.com/iyear/E5SubBot
`
)

var (
	ErrorTimes  map[string]int //错误次数
	BindMaxNum  int
	ErrMaxTimes int
	Notice      string
	Admins      []int64
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {

	}
	viper.SetDefault("errlimit", 5)
	viper.SetDefault("bindmax", 5)

	BindMaxNum = viper.GetInt("bindmax")
	ErrMaxTimes = viper.GetInt("errlimit")
	Notice = viper.GetString("notice")
	Admins = getAdmins()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		BindMaxNum = viper.GetInt("bindmax")
		ErrMaxTimes = viper.GetInt("errlimit")
		Notice = viper.GetString("notice")
		Admins = getAdmins()
	})
}
func getAdmins() []int64 {
	var result []int64
	admins := strings.Split(viper.GetString("admin"), ",")
	for _, v := range admins {
		id, _ := strconv.ParseInt(v, 10, 64)
		result = append(result, id)
	}
	return result
}
