package cryo

import (
	"github.com/machinacanis/cryobot/client"
	"github.com/machinacanis/cryobot/event"
	"github.com/machinacanis/cryobot/log"
)

// HandlerInfo 处理器信息结构体
type HandlerInfo struct {
	Handler     client.EventHandler
	SupportType event.CryoEventType
}

// Processor 处理器结构体
type Processor struct {
	MatchedTypes []event.CryoEventType
	HandlerInfos []HandlerInfo
}

func (h *Processor) Handle(handler interface{}) *Processor {
	switch typedHandler := handler.(type) {
	case func(event.PrivateMessageEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.PrivateMessageEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.PrivateMessageEventType,
		})
	case func(event.GroupMessageEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupMessageEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupMessageEventType,
		})
	case func(event.TempMessageEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.TempMessageEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.TempMessageEventType,
		})
	case func(event.NewFriendRequestEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.NewFriendRequestEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.NewFriendRequestEventType,
		})
	case func(event.NewFriendEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.NewFriendEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.NewFriendEventType,
		})
	case func(event.FriendRecallEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.FriendRecallEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.FriendRecallEventType,
		})
	case func(event.FriendRenameEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.FriendRenameEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.FriendRenameEventType,
		})
	case func(event.FriendPokeEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.FriendPokeEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.FriendPokeEventType,
		})
	case func(event.GroupMemberPermissionUpdatedEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupMemberPermissionUpdatedEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupMemberPermissionUpdatedEventType,
		})
	case func(event.GroupNameUpdatedEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupNameUpdatedEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupNameUpdatedEventType,
		})
	case func(event.GroupMuteEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupMuteEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupMuteEventType,
		})
	case func(event.GroupRecallEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupRecallEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupRecallEventType,
		})
	case func(event.GroupMemberJoinRequestEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupMemberJoinRequestEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupMemberJoinRequestEventType,
		})
	case func(event.GroupMemberIncreaseEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupMemberIncreaseEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupMemberIncreaseEventType,
		})
	case func(event.GroupMemberDecreaseEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupMemberDecreaseEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupMemberDecreaseEventType,
		})
	case func(event.GroupDigestEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupDigestEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupDigestEventType,
		})
	case func(event.GroupReactionEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupReactionEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupReactionEventType,
		})
	case func(event.GroupMemberSpecialTitleUpdated):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupMemberSpecialTitleUpdated); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupMemberSpecialTitleUpdatedEventType,
		})
	case func(event.GroupInviteEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.GroupInviteEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.GroupInviteEventType,
		})
	case func(event.BotConnectedEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.BotConnectedEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.BotConnectedEventType,
		})
	case func(event.BotDisconnectedEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.BotDisconnectedEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.BotDisconnectedEventType,
		})
	case func(event.CustomEvent):
		wrapper := func(e event.CryoEvent) {
			if evt, ok := e.(*event.CustomEvent); ok {
				typedHandler(*evt)
			}
		}
		h.HandlerInfos = append(h.HandlerInfos, HandlerInfo{
			Handler:     wrapper,
			SupportType: event.CustomEventType,
		})
	default:
		log.Warn("注册了不支持的处理器类型！")
	}
	return h
}

// Register 注册所有处理器
func (h *Processor) Register() {
	for _, handlerInfo := range h.HandlerInfos {
		// 只有当处理器的事件类型在匹配类型中时才注册
		for _, matchedType := range h.MatchedTypes {
			if handlerInfo.SupportType == matchedType {
				client.SubscribeSpecific[event.CryoEvent](matchedType, handlerInfo.Handler)
				log.Infof("已注册处理函数到事件类型: %v", matchedType)
			}
		}
	}
}

func (h *Processor) AddEventType(eventType event.CryoEventType) {
	h.MatchedTypes = append(h.MatchedTypes, eventType)
}

func (h *Processor) AddEventTypes(eventTypes ...event.CryoEventType) {
	h.MatchedTypes = append(h.MatchedTypes, eventTypes...)
}

func (h *Processor) SetEventTypes(eventTypes ...event.CryoEventType) {
	h.MatchedTypes = eventTypes
}

func OnType(eventType ...event.CryoEventType) *Processor {
	return &Processor{
		MatchedTypes: eventType,
	}
}
