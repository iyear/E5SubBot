# E5SubBot

![](https://img.shields.io/badge/language-go-blue.svg)
![](https://img.shields.io/badge/license-GPL-lightgrey.svg)

A Simple Telebot for E5 Renewal

Golang + MySQL

DEMO: https://t.me/E5Sub_bot

[交流群组telegram](https://t.me/e5subbot)



## 特性

- 自动续订E5订阅(每三小时调用一次)
- 可管理的简易账户系统
- 极为方便的授权方式

## 原理

E5订阅为开发者订阅，只要调用相关API就有可能续期

调用 [Outlook ReadMail API](https://docs.microsoft.com/zh-cn/graph/api/user-list-messages?view=graph-rest-1.0&tabs=http) 实现玄学的续订方式，不保证续订效果。

## 使用方法

在机器人对话框输入**/bind**，进入授权页面，用E5订阅的账户登录。

授权后会跳转至`http://localhost/e5sub……`

复制整个浏览框内容，在机器人对话框回复 `链接+空格+别名(用于管理账户)`

例如：`http://localhost/e5sub/?code=abcd mye5`，等待机器人绑定后即完成

## 自行部署

#### 二进制文件

在Releases页面下载对应系统的二进制文件，上传至服务器

Windows: 双击`.exe`并保持后台运行

Linux: 

```bash
screen -S e5sub
./e5sub
(Ctrl A+D)
```

#### 编译

下载源码，安装好环境

```shell
go build main.go
```

## 部署配置

在根目录下创建`config.yml`

配置模板

```yaml
bot_token: xxxxx
#不需要socks5代理删去即可
socks5: 127.0.0.1:1080
#auth_url需要自己去Azure注册应用配置
auth_url: https://login.microsoftonline.com/common/oauth2/v2.0/authorize?……
#最大可绑定数
bindmax: 3
#mysql配置
mysql:
  host: 127.0.0.1
  port: 3306
  user: e5sub
  password: e5sub
  database: e5sub
```

## 注意事项

待填写

## 开发过程

由于自己有这方面的需求所以写了一个Bot，Telebot真的很有意思

这是我第二个Go项目(第一个是TG翻译机器人，没几行代码)，原来一直都是小打小闹，也没有正经写过项目。

以前一直没有用过JetBrains IDE，不得不说JB的IDE是真的方便，不过JB+Chrome 8G内存的电脑有点扛不住

Github很早以前就有一个号(记不得了)，但是一直属于只看不用的状态，这次终于也在Github上提交代码了。

Git SQL Golang 基本都是在这四五天里边写边查边学的。目前Git也只会Add Commit Push Pull基本的操作，有一次版本回退折腾了很久……

SQL也只会CURD，代码也是半写半改网上的。

Go在寒假里看过一段时间的《Golang核心编程》+Go官方的教程，但到了现在才开始实践。

------

代码写的很乱也很垃圾，也没有项目的概念，第三方库各种乱调用，文件全是一坨，但实在懒得折腾了。

自行部署可能会有报错崩溃什么的，谅解一下吧。。用DEMO也行

如果有兴趣的大佬欢迎加入群组交流，但有些东西我可能听不太懂。

------

一开始用的sqlite数据库，加载的是<https://github.com/mattn/go-sqlite3> 驱动，结果等我写完了，CGO各种编译错误。

最后实在折腾不起来，只好改用MySQL（VPS上都是一键，Win上装MySQL也折腾了半天）

------

虽然自知写的垃圾，但也是第一个正式的Go Project，还是挺自豪的……

## License

GPLv3 