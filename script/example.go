/*

example.go

这是一个cryobot框架的简单使用示例，编译即可获得一个可运行的简易bot
你可以尝试运行这个bot并连接你的账号，成功连接后，它应该会在你私聊它时复读你的消息

*/

package main

import (
	"github.com/machinacanis/cryobot"
	"github.com/machinacanis/cryobot/client"
	"github.com/machinacanis/cryobot/config"
	"github.com/machinacanis/cryobot/event"
	"github.com/machinacanis/cryobot/log"
	"github.com/sirupsen/logrus"
)

func main() {
	cryo.Init(config.CryoConfig{
		LogLevel:                     logrus.DebugLevel,
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   true,
	})

	client.Subscribe(event.BotConnectedEventType, func(e event.CryoEvent) {
		log.Info("Bot连接事件")
	})

	cryo.OnType(event.BotConnectedEventType, event.PrivateMessageEventType).
		Handle(func(e event.BotConnectedEvent) { log.Info("Bot连接事件") }).
		Handle(func(e event.PrivateMessageEvent) { log.Info("群聊消息") }).
		Register()

	cryo.AutoConnect()
	cryo.Start()
}
