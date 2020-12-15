package main

const jsonData = `{
		"zh_CN": {
			"LC_MESSAGES": {
				"resources.json": [{
					"msgid"       : "helpContent",
					"msgstr"      : ["命令：\n/my 查看已绑定账户信息\n/bind  绑定新账户\n/unbind 解绑账户\n/export 导出账户信息(JSON格式)\n/help 帮助\n/task 手动执行一次任务(管理员)\n/log 获取最近日志文件(管理员)\n源代码：https://github.com/rainerosion/E5SubBot\n原作者：https://github.com/iyear/E5SubBot"]
				},{
					"msgid"       : "welcome",
					"msgstr"      : ["欢迎使用E5SubBot!"]
				},{
					"msgid"       : "chooseAnAccount",
					"msgstr"      : ["选择一个账户查看具体信息\n\n绑定数: "]
				},{
					"msgid"       : "accountInformation",
					"msgstr"      : ["信息\n别名"]
				},{
					"msgid"       : "updateTime",
					"msgstr"      : ["\n最近更新时间: "]
				},{
					"msgid"       : "register",
					"msgstr"      : ["应用注册： [点击直达]"]
				},{
					"msgid"       : "bind1Reply",
					"msgstr"      : ["请回复client_id+空格+client_secret"]
				},{
					"msgid"       : "formatError",
					"msgstr"      : ["错误的格式"]
				},{
					"msgid"       : "signIn",
					"msgstr"      : ["授权账户： [点击直达]"]
				},{
					"msgid"       : "bind2Reply",
					"msgstr"      : ["请回复http://localhost/…… + 空格 + 别名(用于管理)"]
				},{
					"msgid"       : "unBind",
					"msgstr"      : ["选择一个账户将其解绑\n\n当前绑定数: "]
				},{
					"msgid"       : "unBindError",
					"msgstr"      : ["解绑失败!"]
				},{
					"msgid"       : "unBindSuccess",
					"msgstr"      : ["解绑成功!"]
				},{
					"msgid"       : "unbound",
					"msgstr"      : ["你还没有绑定过账户嗷~"]
				},{
					"msgid"       : "json",
					"msgstr"      : ["获取JSON失败\n"]
				},{
					"msgid"       : "temporary",
					"msgstr"      : ["写入临时文件失败~\n"]
				},{
					"msgid"       : "getHelp",
					"msgstr"      : ["发送/help获取帮助嗷"]
				},{
					"msgid"       : "replyBind",
					"msgstr"      : ["请通过回复方式绑定"]
				},{
					"msgid"       : "maximum",
					"msgstr"      : ["已经达到最大可绑定数"]
				},{
					"msgid"       : "binding",
					"msgstr"      : ["正在绑定中……"]
				},{
					"msgid"       : "bindSuccess",
					"msgstr"      : ["绑定成功!"]
				},{
					"msgid"       : "logs",
					"msgstr"      : ["选择一个日志"]
				},{
					"msgid"       : "noPermission",
					"msgstr"      : ["您没有权限执行此操作~"]
				},{
					"msgid"       : "bindFormatError",
					"msgstr"      : ["绑定格式错误"]
				},{
					"msgid"       : "getToken",
					"msgstr"      : ["Token获取成功!"]
				},{
					"msgid"       : "alreadyBind",
					"msgstr"      : ["该应用已经绑定过了，无需重复绑定!"]
				},{
					"msgid"       : "taskError",
					"msgstr"      : ["您的账户: %s 在任务执行时出现了错误!\n错误:"]
				},{
					"msgid"       : "clickToUnBind",
					"msgstr"      : ["点击解绑该账户"]
				},{
					"msgid"       : "unBindByMaxLimit",
					"msgstr"      : ["您的账户因达到错误上限而被自动解绑\n后会有期!\n\n别名: "]
				},{
					"msgid"       : "taskFeedback",
					"msgstr"      : ["任务反馈\n时间: "]
				},{
					"msgid"       : "result",
					"msgstr"      : ["\n结果: "]
				},{
					"msgid"       : "taskFeedbackAdmin",
					"msgstr"      : ["任务反馈(管理员)\n完成时间: "]
				},{
					"msgid"       : "wrongAccount",
					"msgstr"      : ["\n错误账户:\n"]
				},{
					"msgid"       : "clearingAccount",
					"msgstr"      : ["\n清退账户:\n"]
				}
				]
			}
		},
		"en_US": {
			"LC_MESSAGES": {
				"resources.json": [{
					"msgid"       : "helpContent",
					"msgstr"      : ["/my View bound account information\n/bind Bind new account\n/unbind Unbind account\n/export Export account information (JSON format)\n/help help\n/task Manually execute a task (Bot Administrator)\n/log Get the most recent log file (Bot Administrator)\nSource Code：https://github.com/rainerosion/E5SubBot\nOriginal Author：https://github.com/iyear/E5SubBot"]
				},{
					"msgid"       : "welcome",
					"msgstr"      : ["Welcome to E5SubBot!"]
				},{
					"msgid"       : "chooseAnAccount",
					"msgstr"      : ["Select an account to view information.\n\nNumber of bindings: "]
				},{
					"msgid"       : "accountInformation",
					"msgstr"      : ["Account information\nAlias: "]
				},{
					"msgid"       : "updateTime",
					"msgstr"      : ["\nUpdate Time: "]
				},{
					"msgid"       : "register",
					"msgstr"      : ["Application registration： [Click here]"]
				},{
					"msgid"       : "bind1Reply",
					"msgstr"      : ["Please reply client_id+Space+client_secret"]
				},{
					"msgid"       : "formatError",
					"msgstr"      : ["Format error"]
				},{
					"msgid"       : "signIn",
					"msgstr"      : ["Login account [Click here]"]
				},{
					"msgid"       : "bind2Reply",
					"msgstr"      : ["Please reply http://localhost/…… + Space + Alias(for management)"]
				},{
					"msgid"       : "unBind",
					"msgstr"      : ["Select the account to be unbound.\n\nNumber of bindings: "]
				},{
					"msgid"       : "unBindError",
					"msgstr"      : ["Unbind failed!"]
				},{
					"msgid"       : "unBindSuccess",
					"msgstr"      : ["Successfully unbound!"]
				},{
					"msgid"       : "unbound",
					"msgstr"      : ["You haven't tied an account yet."]
				},{
					"msgid"       : "json",
					"msgstr"      : ["Failed to get json\n"]
				},{
					"msgid"       : "temporary",
					"msgstr"      : ["Write to temporary file failed.\n"]
				},{
					"msgid"       : "getHelp",
					"msgstr"      : ["Send /help for help"]
				},{
					"msgid"       : "replyBind",
					"msgstr"      : ["Please bind by reply."]
				},{
					"msgid"       : "maximum",
					"msgstr"      : ["The maximum number of bindings has been reached."]
				},{
					"msgid"       : "binding",
					"msgstr"      : ["Binding……"]
				},{
					"msgid"       : "bindSuccess",
					"msgstr"      : ["Binding succeeded!"]
				},{
					"msgid"       : "logs",
					"msgstr"      : ["Select a log."]
				},{
					"msgid"       : "noPermission",
					"msgstr"      : ["You do not have permission."]
				},{
					"msgid"       : "bindFormatError",
					"msgstr"      : ["Format error"]
				},{
					"msgid"       : "getToken",
					"msgstr"      : ["Token obtained successfully!"]
				},{
					"msgid"       : "alreadyBind",
					"msgstr"      : ["Account already exists!"]
				},{
					"msgid"       : "taskError",
					"msgstr"      : ["The account %s an error.\nError details:"]
				},{
					"msgid"       : "clickToUnBind",
					"msgstr"      : ["Click to unbundle this account"]
				},{
					"msgid"       : "unBindByMaxLimit",
					"msgstr"      : ["Your account has been automatically unbundled due to an error limit being reached\nAccount alias: "]
				},{
					"msgid"       : "taskFeedback",
					"msgstr"      : ["Task feedback\nTime: "]
				},{
					"msgid"       : "result",
					"msgstr"      : ["\nResult: "]
				},{
					"msgid"       : "taskFeedbackAdmin",
					"msgstr"      : ["Task feedback(administrator)\nExecution time: "]
				},{
					"msgid"       : "wrongAccount",
					"msgstr"      : ["\nWrong account:\n"]
				},{
					"msgid"       : "clearingAccount",
					"msgstr"      : ["\nClearing account:\n"]
				}
				]
			}
		}
	}`
