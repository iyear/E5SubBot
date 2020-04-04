# E5SubBot

![](https://img.shields.io/github/go-mod/go-version/iyear/E5SubBot)
![](https://img.shields.io/badge/license-GPL-lightgrey.svg)
![](https://img.shields.io/github/v/release/iyear/E5SubBot?color=green)

A Simple Telebot for E5 Renewal

Golang + MySQL

DEMO: https://t.me/E5Sub_bot(长期运行，所有新功能会在DEMO测试)

[交流群组telegram](https://t.me/e5subbot)

## 预览
- [绑定过程](https://raw.githubusercontent.com/iyear/E5SubBot/master/pics/bind.JPG)
- [查看"我的"](https://raw.githubusercontent.com/iyear/E5SubBot/master/pics/my.JPG)
- [任务反馈](https://raw.githubusercontent.com/iyear/E5SubBot/master/pics/task.JPG)

## 特性

- 自动续订E5订阅(可自定义的调用频率)
- 可管理的简易账户系统
- 完善的任务执行反馈
- 极为方便的授权方式


## 原理

E5订阅为开发者订阅，只要调用相关API就有可能续期

调用 [Outlook ReadMail API](https://docs.microsoft.com/zh-cn/graph/api/user-list-messages?view=graph-rest-1.0&tabs=http) 实现玄学的续订方式，不保证续订效果。

## 使用方法

1. 在机器人对话框输入 **/bind**
2. 注册应用，使用E5主账号或同域账号登录，跳转页面获得client_secret。**点击回到快速启动**,获得client_id
3. 复制client_secret和client_id，以 `client_id client_secret`格式回复
4. 获得授权链接，使用E5主账号或同域账号登录
5. 授权后会跳转至`http://localhost/e5sub……`
6. 复制整个浏览框内容，在机器人对话框回复 `链接+空格+别名(用于管理账户)`

例如：`http://localhost/e5sub/?code=abcd MyE5`，等待机器人绑定后即完成

## 自行部署
需要MySQL>=5.5版本(开发本地环境是5.5，高的低的没测试过，应该也可以)

Bot创建教程:[Google](https://www.google.com/search?q=telegram+Bot%E5%88%9B%E5%BB%BA%E6%95%99%E7%A8%8B)
#### 二进制文件

在[Releases](https://github.com/iyear/E5SubBot/releases)页面下载对应系统的二进制文件，上传至服务器

Windows: 在cmd中启动 `E5SubBot.exe`

Linux: 

```bash
screen -S e5sub
chmod 773 E5SubBot
./E5SubBot
(Ctrl A+D)
```
#### 编译

下载源码，安装好环境

```shell
go build main.go
```

## 部署配置

在根目录下创建`config.yml`，编码为UTF-8

配置模板

```yaml
#bindmax,notice,admin可热更新，直接更新config.yml保存即可
#更换为自己的BotToken
bot_token: xxxxx
#不需要socks5代理删去即可
socks5: 127.0.0.1:1080
#公告，合并至/help
notice: "第一行\n第二行"
#管理员tgid，前往https://t.me/userinfobot获取，用,隔开
#管理员权限: 手动调用任务，获得任务总反馈
admin: 66666,77777,88888
#API调用频率，使用cron表达式
cron: "1 */3 * * *"
#最大可绑定数
bindmax: 3
#mysql配置，请提前创建数据库
mysql:
  host: 127.0.0.1
  port: 3306
  user: e5sub
  password: e5sub
  database: e5sub
```

## 注意事项
> 更新时间与北京时间不符

更改服务器时区为Asia/Shanghai，然后使用/task手动执行一次任务


## License

GPLv3 