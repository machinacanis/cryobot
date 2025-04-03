/*

example.go

这是一个cryobot框架的简单使用示例，编译即可获得一个可运行的简易bot
你可以尝试运行这个bot并连接你的账号，成功连接后，它应该会在你私聊它时复读你的消息

*/

package main

import (
	"github.com/machinacanis/cryobot"
	"github.com/machinacanis/cryobot/config"
	"github.com/sirupsen/logrus"
)

func main() {
	cryo.Init(config.CryoConfig{
		LogLevel:                     logrus.InfoLevel,
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   true,
	})

	cryo.AutoConnect()
	cryo.Start()
}
