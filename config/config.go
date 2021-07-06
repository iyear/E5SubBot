package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	BotToken      string
	Socks5        string
	BindMaxNum    int
	MaxGoroutines int
	MaxErrTimes   int
	Cron          string
	Notice        string
	Admins        []int64
	Mysql         MysqlConfig
)

type MysqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
	Table    string
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		zap.S().Errorw("failed to read config", "error", err)
	}
	BotToken = viper.GetString("bot_token")
	Cron = viper.GetString("cron")
	Socks5 = viper.GetString("socks5")

	viper.SetDefault("errlimit", 5)
	viper.SetDefault("bindmax", 5)
	viper.SetDefault("goroutine", 10)

	BindMaxNum = viper.GetInt("bindmax")
	MaxErrTimes = viper.GetInt("errlimit")
	Notice = viper.GetString("notice")

	MaxGoroutines = viper.GetInt("goroutine")
	Admins = getAdmins()

	Mysql = MysqlConfig{
		Host:     viper.GetString("mysql.host"),
		Port:     viper.GetInt("mysql.port"),
		User:     viper.GetString("mysql.user"),
		Password: viper.GetString("mysql.password"),
		DB:       viper.GetString("mysql.database"),
		Table:    viper.GetString("mysql.table"),
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		MaxGoroutines = viper.GetInt("goroutine")
		BindMaxNum = viper.GetInt("bindmax")
		MaxErrTimes = viper.GetInt("errlimit")
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
