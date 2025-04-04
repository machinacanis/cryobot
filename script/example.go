package main

import (
	cryo "github.com/machinacanis/cryobot"
	"github.com/sirupsen/logrus"
)

func main() {
	cryo.Init(cryo.Config{
		LogLevel:                     logrus.DebugLevel,
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   true,
	})
	cryo.AutoConnect()
	cryo.Start()
}
