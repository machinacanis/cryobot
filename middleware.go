package cryobot

// setConnectPrintMiddleware 内置的连接打印中间件
func setConnectPrintMiddleware() {
	if conf.EnableConnectPrintMiddleware {
		AddMiddleware(BotConnectedEventType, func(e CryoEvent) CryoEvent {
			if typedEvent, ok := e.(BotConnectedEvent); ok {
				Infof("%s[Cryo] %s：%s (%d) 已成功连接", lavender, typedEvent.BotNickname, typedEvent.BotId, typedEvent.BotUin)
			}
			return e
		})
		AddMiddleware(BotDisconnectedEventType, func(e CryoEvent) CryoEvent {
			if typedEvent, ok := e.(BotDisconnectedEvent); ok {
				Infof("%s[Cryo] %s：%s (%d) 已断开连接", lavender, typedEvent.BotNickname, typedEvent.BotId, typedEvent.BotUin)
			}
			return e
		})
	}
}

// setMessagePrintMiddleware 内置的消息打印中间件
func setMessagePrintMiddleware() {
	if conf.EnableMessagePrintMiddleware {
		AddMiddleware(PrivateMessageEventType, func(e CryoEvent) CryoEvent {
			if typedEvent, ok := e.(PrivateMessageEvent); ok {
				Infof("[%s] [私聊] From %s(%d) - %s", typedEvent.BotNickname, typedEvent.SenderNickname, typedEvent.SenderUin, typedEvent.MessageElements.ToString())
			}
			return e
		})
		AddMiddleware(GroupMessageEventType, func(e CryoEvent) CryoEvent {
			if typedEvent, ok := e.(GroupMessageEvent); ok {
				Infof("[%s] [%s(%d)] From %s(%d) - %s", typedEvent.BotNickname, typedEvent.GroupName, typedEvent.GroupUin, typedEvent.SenderNickname, typedEvent.SenderUin, typedEvent.MessageElements.ToString())
			}
			return e
		})
	}
}

// setEventDebugMiddleware 内置的事件调试中间件
func setEventDebugMiddleware() {
	if conf.EnableEventDebugMiddleware {
		AddGlobalMiddleware(func(e CryoEvent) CryoEvent {
			Debug(e.ToJsonString()) // 输出json
			return e
		})
	}
}
