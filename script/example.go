/*

example.go

这是一个cryobot框架的简单使用示例，编译即可获得一个可运行的简易bot
你可以尝试运行这个bot并连接你的账号，成功连接后，它应该会在你私聊它时复读你的消息

*/

package main

import (
	"github.com/machinacanis/cryobot/client"
	"github.com/machinacanis/cryobot/event"
	"github.com/machinacanis/cryobot/log"
	"github.com/sirupsen/logrus"
)

func main() {
	log.InitTextLogger(logrus.DebugLevel)
	log.Info("欢迎使用cryobot！")
	client.SubscribeSpecific(event.PrivateMessageEventType, func(event event.PrivateMessageEvent) {
		// Handle private message
		log.Info(event.ToJsonString())
	})

	client.SubscribeSpecific(event.GroupMessageEventType, func(event event.GroupMessageEvent) {
		// Handle group message
		log.Info(event.ToJsonString())
	})

	client.ConnectAll()
	select {}
}
