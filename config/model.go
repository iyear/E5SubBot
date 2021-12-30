package config

var (
	BotToken      string
	Socks5        string
	BindMaxNum    int
	MaxGoroutines int
	MaxErrTimes   int
	Cron          string
	Notice        string
	Admins        []int64
	DB            string
	Table         string
	Mysql         mysqlConfig
	Sqlite        sqliteConfig
)

type sqliteConfig struct {
	DB string `json:"db,omitempty"`
}
type mysqlConfig struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	DB       string `json:"db,omitempty"`
}
