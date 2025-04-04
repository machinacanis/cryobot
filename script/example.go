package main

import (
	cryo "github.com/machinacanis/cryobot"
	"github.com/sirupsen/logrus"
)

func main() {
	bot := cryo.NewBot()
	bot.Init(cryo.Config{
		LogLevel:                     logrus.DebugLevel,
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   true,
	})

	bot.OnType(cryo.PrivateMessageEventType).
		Handle(func(e cryo.PrivateMessageEvent) {
			cryo.Info("接收到了" + e.SenderNickname + "的私聊消息！")
		}).
		Handle(func(e cryo.GroupMessageEvent) {
			cryo.Info("接收到了" + e.GroupName + "的群消息！")
		}).
		Register()
	bot.AutoConnect()
	bot.Start()
}
