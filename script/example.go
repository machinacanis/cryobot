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

	bot.OnMessage().
		HandleMessage(func(e cryo.MessageEvent) {
			if e.GroupUin == 941419619 {
				bot.Reply(e, "你说的对")
			}
		}).
		Register()

	bot.AutoConnect()
	bot.Start()
}
