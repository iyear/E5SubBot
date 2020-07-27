# E5SubBot

![](https://img.shields.io/github/go-mod/go-version/iyear/E5SubBot?style=flat-square)
![](https://img.shields.io/badge/license-GPL-lightgrey.svg?style=flat-square)
![](https://img.shields.io/github/v/release/iyear/E5SubBot?color=green&style=flat-square)

English | [简体中文](https://github.com/iyear/E5SubBot/blob/master/README_zhCN.md)

A Simple Telebot for E5 Renewal

`Golang` + `MySQL`

DEMO: https://t.me/E5Sub_bot (all new functions will be tested in DEMO)

Communication: [Telegram Group](https://t.me/e5subbot)

## Preview
<center class="half">
    <img src="https://raw.githubusercontent.com/iyear/E5SubBot/master/pics/bind.JPG" width="200"/><img src="https://raw.githubusercontent.com/iyear/E5SubBot/master/pics/my.JPG" width="200"/><img src="https://raw.githubusercontent.com/iyear/E5SubBot/master/pics/task.JPG" width="200"/>
</center>

## Feature

- Automatically Renew E5 Subscription(Customizable Frequency)
- Manageable Simple Account System
- Available Task Execution Feedback
- Convenient Authorization


## Principle

E5 subscription is a subscription for developers, as long as the related API is called, it may be renewed

Calling [Outlook ReadMail API](https://docs.microsoft.com/en-us/graph/api/user-list-messages?view=graph-rest-1.0&tabs=http) to renew, does not guarantee the renewal effect.

## Usage

1. Type `/bind` in the robot dialog
2. Click the link sent by the robot and register the Microsoft application, log in with the E5 master account or the same domain account, and obtain `client_secret`. **Click to go back to Quick Start**, get `client_id`
3. Copy `client_secret` and `client_id` and reply to bot in the format of `client_id(space)client_secret`
(Pay attention to spaces)
4. Click on the authorization link sent by the robot and log in with the `E5` master account or the same domain account
5. After authorization, it will jump to `http://localhost/e5sub……` (will prompt webpage error, just copy the link)
6. Copy the link, and reply `link(space)alias (used to manage accounts)` in the robot dialog
For example: `http://localhost/e5sub/?code=abcd MyE5`, wait for the robot to bind and then complete

## Deploy Your Own Bot 

Bot creation tutorial : [Microsoft](https://docs.microsoft.com/en-us/azure/bot-service/bot-service-channel-connect-telegram?view=azure-bot-service-4.0)

### Docker Deployment
Thanks to [@kzw200015](https://github.com/kzw200015) for providing help in `Dockerfile` and `Docker`

If it fails to start for the first time, use `docker-compose restart` to restart
```bash
mkdir ./e5bot && wget --no-check-certificate -O ./e5bot/config.yml https://raw.githubusercontent.com/iyear/E5SubBot/master/config.yml.example
vi ./e5bot/config.yml
wget --no-check-certificate https://raw.githubusercontent.com/iyear/E5SubBot/master/docker-compose.yml
docker-compose up -d
```
### Binary Deployment

Download the binary files of the corresponding system on the [Releases](https://github.com/iyear/E5SubBot/releases) page and upload it to the server

Windows: Start `E5SubBot.exe` in `cmd`

Linux: 

```bash
screen -S e5sub
chmod 773 E5SubBot
./E5SubBot
(Ctrl A+D)
```
### Compile

Download the source code and install the GO environment

```shell
go build
```

## Configuration

Create `config.yml` in the same directory, encoded as `UTF-8`

Configuration Template:

```yaml
bot_token: YOUR_BOT_TOKEN
socks5: 127.0.0.1:1080
notice: "first line \n second line"
admin: 66666,77777,88888
errlimit: 5
cron: "1 */3 * * *"
bindmax: 3
mysql:
  host: 127.0.0.1
  port: 3306
  user: e5sub
  password: e5sub
  database: e5sub
```

`bindmax`, `notice`, `admin`, `errlimit` can be hot updated, just update `config.yml` to save.
|  Configuration   | Explanation|
|  ----  | ----  |
| bot_token  | Change to your own `BotToken` |
| socks5  | `Socks5` proxy,if you do not need ,you should delete it. For example: `127.0.0.1:1080` |
|notice|Announcement. Merged into `/help`|
|admin|The administrator's `tgid`, go to https://t.me/userinfobot to get it, separated by `,`; Administrator permissions: manually call the task, get the total feedback of the task|
|errlimit|The maximum number of errors for a single account, automatically unbind the single account and send a notification when it is full, without limiting the number of errors, change the value to a negative number `(-1)`; all errors will be cleared after the bot restarts|
|cron|API call frequency, using `cron` expression|
|bindmax|Maximum number of bindable|
|mysql|Mysql configuration, please create database in advance|

### Command
```
/my View bound account information
/bind Bind new account
/unbind Unbind account
/export Export account information (JSON format)
/help help
/task Manually execute a task (Bot Administrator)
/log Get the most recent log file (Bot Administrator)
```
## Others
> Feedback time is not as expected

Change the server time zone, use `/task` to manually perform a task to refresh time.

> ERROR:Can't create more than max_prepared_stmt_count statements (current value: 16382).

Failure to close `db` leads to triggering `mysql` concurrency limit, please update to `v0.1.9`.

> Long running crash

Suspected memory leak. Not yet resolved, please run the daemon or restart Bot regularly.

> Unable to create application via bot

https://t.me/e5subbot/5201

## Third-Party
- [Telebot](https://gopkg.in/tucnak/telebot)
- [Mysql_driver](https://github.com/go-sql-driver/mysql)
- [Gjson](https://github.com/tidwall/gjson)
- [Cron](https://github.com/robfig/cron/)
- [Viper](https://github.com/spf13/viper)
- [Goreleaser](https://https://github.com/goreleaser/goreleaser)

## Contributing
- Provide documentation in other languages
- Provide help for code operation
- Suggests user interaction
- ……
## More Functions
If you still want to support new features, please use FeatHub to vote. We will consider the voting results and other factors to determine the development priority.  

[![Feature Requests](https://cloud.githubusercontent.com/assets/390379/10127973/045b3a96-6560-11e5-9b20-31a2032956b2.png)](http://feathub.com/NervJS/taro)  

[![Feature Requests](https://feathub.com/iyear/E5SubBot?format=svg)](https://feathub.com/iyear/E5SubBot)  

## License

GPLv3 
