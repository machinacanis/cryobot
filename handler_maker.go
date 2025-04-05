package cryobot

var messageEventTypes = []CryoEventType{
	PrivateMessageEventType,
	GroupMessageEventType,
	TempMessageEventType,
}

// On 创建一个空的事件处理器
func (b *Bot) On() *Handler {
	return &Handler{}
}

// OnType 创建一个可以匹配类型的事件处理器
func (b *Bot) OnType(eventType ...CryoEventType) *Handler {
	return &Handler{
		MatchingTypes: eventType, // 事件类型
	}
}

// OnMessage 创建一个消息事件处理器
func (b *Bot) OnMessage(eventType ...CryoEventType) *Handler {
	if len(eventType) == 0 {
		eventType = messageEventTypes
	} else if len(eventType) > 0 {
		// 如果传入的事件类型不在消息事件类型列表中，则返回默认的消息事件处理器
		for _, et := range eventType {
			if !Contains(messageEventTypes, et) {
				eventType = messageEventTypes
				break
			}
		}
	}
	return &Handler{
		MatchingTypes: eventType, // 事件类型
	}
}

func (b *Bot) OnFullmatch(text ...string) *Handler {
	return &Handler{}
}
