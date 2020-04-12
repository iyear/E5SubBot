# E5SubBot

![](https://img.shields.io/github/go-mod/go-version/iyear/E5SubBot)
![](https://img.shields.io/badge/license-GPL-lightgrey.svg)
![](https://img.shields.io/github/v/release/iyear/E5SubBot?color=green)

A Simple Telebot for E5 Renewal

Golang + MySQL

DEMO: https://t.me/E5Sub_bot (长期运行，所有新功能会在DEMO测试)

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
### 命令
```
/my 查看已绑定账户信息
/bind  绑定新账户
/unbind 解绑账户
/export 导出账户信息(json格式)
/help 帮助
/task 手动执行一次任务(管理员)
/log 获取最近日志文件(管理员)
```
## 自行部署
需要MySQL>=5.5版本(开发本地环境是5.5，高的低的没测试过，应该也可以)
Bot创建教程:[Google](https://www.google.com/search?q=telegram+Bot%E5%88%9B%E5%BB%BA%E6%95%99%E7%A8%8B)

### Docker部署
感谢 [@kzw200015](https://github.com/kzw200015) 提供`Dockerfile`以及`Docker`方面的帮助

第一次启动不行，使用 `docker-compose restart`重启一次
```bash
mkdir ./e5bot && wget --no-check-certificate -O ./e5bot/config.yml https://raw.githubusercontent.com/iyear/E5SubBot/master/config.yml.example
vi ./e5bot/config.yml
wget --no-check-certificate https://raw.githubusercontent.com/iyear/E5SubBot/master/docker-compose.yml
docker-compose up -d
```
### 二进制文件

在[Releases](https://github.com/iyear/E5SubBot/releases)页面下载对应系统的二进制文件，上传至服务器

Windows: 在cmd中启动 `E5SubBot.exe`

Linux: 

```bash
screen -S e5sub
chmod 773 E5SubBot
./E5SubBot
(Ctrl A+D)
```
### 编译

下载源码，安装好环境

```shell
go build
```

## 部署配置

在根目录下创建`config.yml`，编码为UTF-8

配置模板

```yaml
#bindmax,notice,admin,errlimit可热更新，直接更新config.yml保存即可
#更换为自己的BotToken
bot_token: xxxxx
#不需要socks5代理删去即可
socks5: 127.0.0.1:1080
#公告，合并至/help
notice: "第一行\n第二行"
#管理员tgid，前往https://t.me/userinfobot获取，用,隔开
#管理员权限: 手动调用任务，获得任务总反馈
admin: 66666,77777,88888
#任务最大出错次数，满后自动解绑账户并发送通知，无限次数将值改为负数(-1)即可
#以ms账户为单位，不会解绑所有账户(只解绑错误账户)
#主要为了减少资源浪费.bot重启后会清零所有错误次数
errlimit: 5
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
### 命令
>   /my 查看已绑定账户信息
 	/bind  绑定新账户
 	/unbind 解绑账户
 	/export 导出账户信息(JSON格式)
 	/help 帮助
 	/task 手动执行一次任务(管理员)
 	/log 获取最近日志文件(管理员)
## 注意事项
> 更新时间与北京时间不符

更改服务器时区为Asia/Shanghai，然后使用/task手动执行一次任务刷新时间

> 绑定格式错误

不要带"+"号

## 小总结
#### 得到了什么？
- git的基本操作:add,commit,pull,push.但是对代码合并与冲突解决没有得到足够的实践和深入  
- github版本库方面的使用(issue,pull request...),还有一些没玩过
- 体验了Docker Hub 的自动构建
- sql的CRUD，但都只是浮于表面，有时间再去琢磨
- telegram bot的golang基本框架和一些有意思的玩法
- 一些著名第三方库(viper,gjson)的基本用法
- docker以及docker-compose的基本用法
- 一些工具:gox;goreleaser的基本用法
- ……
#### 还缺点什么？

- 对`CGO`知难而退，编译各种出错。下个项目如果用到`sqlite`一定解决`CGO`交叉编译
- `DockerFile` 还不会写
- 对项目的概念和意识不太行，目录太乱，随缘写全局变量
- 看`telebot`的文档还好，看`stackflow` 只能看看代码。搜索还是习惯性带中文导致`stackflow`很难出现，降低了解决问题的效率

#### 最后

从2020.3.28开发至2020.4.12，一共经历14天，利用课余时间最终完成。

就要开学了，因为马上步入高三，学习紧迫，遂不再进行功能开发。

项目差不多就停更了，最多也就是每个星期看看issue、tg群组，修一修bug

**DEMO依旧能保持服务水平，不会因为个人问题在这段时间停止运行**

如果一年后还有人在用，一定会继续

## Third-Party
- [telebot](https://gopkg.in/tucnak/telebot)
- [mysql_driver](https://github.com/go-sql-driver/mysql)
- [gjson](https://github.com/tidwall/gjson)
- [cron](https://github.com/robfig/cron/)
- [viper](https://github.com/spf13/viper)
- [goreleaser](https://https://github.com/goreleaser/goreleaser)

## License

GPLv3 
