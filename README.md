<img src="https://github.com/iyear/E5SubBot/raw/master/pics/office.png" alt="logo" width="130" height="130" align="left" />

<h1>E5SubBot</h1>

> A Simple Telebot for E5 Renewal

<br/>

![](https://img.shields.io/github/go-mod/go-version/iyear/E5SubBot?style=flat-square)
![](https://img.shields.io/badge/license-GPL-lightgrey.svg?style=flat-square)
![](https://img.shields.io/github/v/release/iyear/E5SubBot?color=red&style=flat-square)
![](https://img.shields.io/github/last-commit/iyear/E5SubBot?style=flat-square)
![](https://img.shields.io/github/downloads/iyear/E5SubBot/total?style=flat-square)

![](https://img.shields.io/github/workflow/status/iyear/E5SubBot/Docker%20Build?label=docker%20build&style=flat-square)
![](https://img.shields.io/docker/v/iyear/e5subbot?label=docker%20tag&style=flat-square)
![](https://img.shields.io/docker/image-size/iyear/e5subbot?style=flat-square&label=docker%20image%20size)

English | [简体中文](https://github.com/iyear/E5SubBot/blob/master/README_zhCN.md) | [Telegram Chat](https://t.me/e5subbot)

DEMO: https://t.me/E5Sub_bot

## Feature

- Automatically Renew E5 Subscription(Customizable Frequency)
- Manageable Simple Account System
- Available Task Execution Feedback
- Convenient Authorization
- Use concurrency to speed up

## Principle

E5 subscription is a subscription for developers, as long as the related API is called, it may be renewed

Calling [Outlook ReadMail API](https://docs.microsoft.com/en-us/graph/api/user-list-messages?view=graph-rest-1.0&tabs=http)
to renew, does not guarantee the renewal effect.

## Usage

1. Type `/bind` in the robot dialog
2. Click the link sent by the robot and register the Microsoft application, log in with the E5 master account or the
   same domain account, and obtain `client_secret`. **Click to go back to Quick Start**, get `client_id`
3. Copy `client_secret` and `client_id` and reply to bot in the format of `client_id(space)client_secret`
   (Pay attention to spaces)
4. Click on the authorization link sent by the robot and log in with the `E5` master account or the same domain account
5. After authorization, it will jump to `http://localhost/e5sub……` (will prompt webpage error, just copy the link)
6. Copy the link, and reply `link(space)alias (used to manage accounts)` in the robot dialog For
   example: `http://localhost/e5sub/?code=abcd MyE5`, wait for the robot to bind and then complete

## Deploy Your Own Bot

Bot creation
tutorial : [Microsoft](https://docs.microsoft.com/en-us/azure/bot-service/bot-service-channel-connect-telegram?view=azure-bot-service-4.0)

### Docker(Recommended)

`Docker` Deployment used `sqlite` as database

Support `amd64` `386` `arm64` `arm/v6` `arm/v7` arch

```shell
#launch,you can set the time zone you want
docker run --name e5sub -e TZ="Asia/Shanghai" --restart=always -d iyear/e5subbot:latest

#view logs
docker logs -f e5sub

#set config
docker cp PATH/TO/config.yml e5sub:/config.yml
docker restart e5sub

#import db
docker cp PATH/TO/DATA.db e5sub:/data.db
docker restart e5sub

#backup db
docker cp e5sub:/data.db .

#backup config
docker cp e5sub:/config.yml .
```

### Binary Deployment

Download the binary files of the corresponding system on the [Releases](https://github.com/iyear/E5SubBot/releases) page
and upload it to the server

Windows: Start `E5SubBot.exe`

Linux:

```bash
screen -S e5sub
chmod +x E5SubBot
./E5SubBot
(Ctrl A+D)
```

### Compile

Download the source code and install the GO environment

```shell
git clone https://github.com/iyear/E5SubBot.git && cd E5SubBot && go build
```

## Configuration

Create `config.yml` in the same directory, encoded as `UTF-8`

Configuration Template:

```yaml
bot_token: YOUR_BOT_TOKEN
# socks5: 127.0.0.1:1080
bindmax: 999
goroutine: 20
admin: 111,222,333
errlimit: 999
notice: |-
   aaa
   bbb
   ccc
cron: "1 */1 * * *"
db: sqlite
table: users
# mysql:
#    host: 127.0.0.1
#    port: 3306
#    user: root
#    password: pwd
#    database: e5sub
sqlite:
   db: data.db
```

`bindmax`, `notice`, `admin`,`goroutine`, `errlimit` can be hot updated, just update `config.yml` to save.

|  Configuration   | Explanation|Default|
|  ----  | ----  |----|
| bot_token  | Change to your own `BotToken` |-|
| socks5  | `Socks5` proxy,if you do not need ,you should delete it. For example: `127.0.0.1:1080` |-|
|notice|Announcement. Merged into `/help`|-|
|admin|The administrator's `tgid`, go to https://t.me/userinfobot to get it, separated by `,`; Administrator permissions: manually call the task, get the total feedback of the task|-|
|goroutine|Concurrent number, don’t be too big|10|
|errlimit|The maximum number of errors for a single account, automatically unbind the single account and send a notification when it is full, without limiting the number of errors, change the value to a negative number `(-1)`; all errors will be cleared after the bot restarts|5|
|cron|API call frequency, using `cron` expression|-|
|bindmax|Maximum number of bindable|5|
|db|`mysql` or `sqlite` , Indicates the database type used and sets the corresponding configuration|-|
|table|Table name (set table to `users` when upgrading the old version; otherwise, the data table cannot be read)|-|
|mysql|To configure `mysql`, create a database in advance|-|
|sqlite|`sqlite` configuration|-|

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

## Contributing

- Provide documentation in other languages
- Provide help for code operation
- Suggests user interaction
- ……

## More Functions

If you still want to support new features, please initiate an issue.

## License

GPLv3 
