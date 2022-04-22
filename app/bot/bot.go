package bot

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/iyear/E5SubBot/app/bot/internal/config"
	"github.com/iyear/E5SubBot/app/bot/internal/logger"
	"github.com/iyear/E5SubBot/app/bot/internal/middleware"
	"github.com/iyear/E5SubBot/app/bot/internal/model"
	"github.com/iyear/E5SubBot/app/bot/internal/template"
	"github.com/iyear/E5SubBot/pkg/conf"
	"github.com/iyear/E5SubBot/pkg/db"
	"github.com/iyear/E5SubBot/pkg/models"
	"github.com/iyear/sqlite"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	tele "gopkg.in/telebot.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path"
	"time"
)

func Run(cfg string, tmplCfg string, dataPath string) {
	color.Blue(conf.LOGO)
	log := logger.Init()

	if err := config.Init(cfg); err != nil {
		log.Fatalw("init config failed", "err", err)
	}
	color.Blue("read config succ...\n")

	if err := template.Init(tmplCfg); err != nil {
		log.Fatalw("init template config failed", "err", err)
	}
	color.Blue("read template succ...\n")

	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		log.Fatalw("create data dir failed", "err", err)
	}

	dial, err := getDialector(dataPath)
	rl, err := db.InitRelational(dial, &models.Client{})
	if err != nil {
		log.Fatalw("init db failed", "err", err)
	}
	color.Blue("init db succ...\n")

	kv, err := db.InitKV(path.Join(dataPath, conf.DBBolt), conf.BucketLanguage)
	if err != nil {
		log.Fatalw("init kv database failed", "err", err, "path", dataPath)
	}
	color.Blue("init kv db succ...\n")

	settings := tele.Settings{
		Token:     config.C.Bot.Token,
		Poller:    &tele.LongPoller{Timeout: 5 * time.Second},
		Client:    getClient(),
		OnError:   middleware.OnError(),
		ParseMode: tele.ModeMarkdown,
	}

	bot, err := tele.NewBot(settings)
	if err != nil {
		log.Fatalw("create bot failed", "err", err)
	}
	color.Blue("create bot succ...\n")

	color.Blue("Bot: %s", bot.Me.Username)

	bot.Use(middleware.SetScope(&model.Scope{
		DB: &model.DB{
			DB:    rl,
			KV:    kv,
			Cache: cache.New(cache.NoExpiration, time.Minute),
		},
		Log: log.Named("handler"),
	}), middleware.AutoResponder())

	makeHandlers(bot)

	bot.Start()
}

func getDialector(dataPath string) (gorm.Dialector, error) {
	switch config.C.DB.Driver {
	case conf.RelationalMySQL:
		c := config.C.DB.MySQL
		return mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User,
			c.Password,
			c.Host,
			c.Port,
			c.Database,
		)), nil
	case conf.RelationalSQLite:
		return sqlite.Open(path.Join(dataPath, conf.DBSQLite)), nil
	}
	return nil, fmt.Errorf("unknown db driver: %s", config.C.DB.Driver)
}

func getClient() *http.Client {
	c := config.C.Bot.Socks5
	if !c.Enable {
		return http.DefaultClient
	}

	host := c.Host
	port := c.Port
	user := c.User
	password := c.Password

	dialer, err := proxy.SOCKS5("tcp",
		fmt.Sprintf("%s:%d", host, port),
		&proxy.Auth{User: user, Password: password},
		proxy.Direct)

	if err != nil {
		zap.S().Fatalw("failed to get dialer",
			"error", err,
			"host", host,
			"port", port,
			"user", user,
			"password", password)
	}
	return &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
}
