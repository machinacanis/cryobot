package cryobot

type Subscription struct {
	HandlerId   string
	HandlerFunc func(CryoEvent)
	HandlerType CryoEventType
}

// TypedWrapper 带泛型的事件处理函数包装器
func TypedWrapper[T CryoEvent](handler func(T)) func(CryoEvent) {
	return func(e CryoEvent) {
		if evt, ok := e.(T); ok {
			handler(evt)
		}
	}
}

// Handler cryobot的事件处理器
type Handler struct {
	Tags               []string        // 事件处理器的标签，这些标签会被带入这个事件处理器生成的订阅中
	Subscriptions      []Subscription  // 将被用于订阅的事件处理函数列表
	Middlewares        []Middleware    // 中间件订阅列表
	MessageMiddlewares []Middleware    // 消息中间件订阅列表
	MatchingTypes      []CryoEventType // 支持处理的事件类型
}

// AddTags 用于向事件处理器添加标签
func (h *Handler) AddTags(tags ...string) *Handler {
	// 将标签添加到事件处理器，如果已经有重复的标签，则不添加
	for _, tag := range tags {
		if !Contains(h.Tags, tag) {
			h.Tags = append(h.Tags, tag)
		}
	}
	return h
}

// GetTags 返回事件处理器的标签
func (h *Handler) GetTags() []string {
	return h.Tags
}

// SetTags 用于直接覆盖设置事件处理器的标签
func (h *Handler) SetTags(tags ...string) *Handler {
	h.Tags = tags
	return h
}

func (h *Handler) AddMatchingTypes(types ...CryoEventType) *Handler {
	// 将匹配的事件类型添加到事件处理器
	for _, et := range types {
		if !Contains(h.MatchingTypes, et) {
			h.MatchingTypes = append(h.MatchingTypes, et)
		}
	}
	return h
}

// GetMatchingTypes 返回事件处理器的匹配事件类型
func (h *Handler) GetMatchingTypes() []CryoEventType {
	return h.MatchingTypes
}

// SetMatchingTypes 用于直接覆盖设置事件处理器的匹配事件类型
func (h *Handler) SetMatchingTypes(types ...CryoEventType) *Handler {
	h.MatchingTypes = types
	return h
}

// AddMiddlewares 用于向事件处理器添加中间件
func (h *Handler) AddMiddlewares(middlewares ...Middleware) *Handler {
	// 将中间件添加到事件处理器
	h.Middlewares = append(h.Middlewares, middlewares...)
	return h
}

func (h *Handler) ClearMiddlewares() *Handler {
	// 清空事件处理器的中间件
	h.Middlewares = []Middleware{}
	return h
}

// AddMessageMiddlewares 用于向事件处理器添加消息中间件
func (h *Handler) AddMessageMiddlewares(middlewares ...Middleware) *Handler {
	// 将消息中间件添加到事件处理器
	h.MessageMiddlewares = append(h.MessageMiddlewares, middlewares...)
	return h
}

// ClearMessageMiddlewares 清空事件处理器的消息中间件
func (h *Handler) ClearMessageMiddlewares() *Handler {
	// 清空事件处理器的消息中间件
	h.MessageMiddlewares = []Middleware{}
	return h
}

// Handle 用于向事件处理器添加处理函数
func (h *Handler) Handle(handler interface{}) *Handler {
	switch typedHandler := handler.(type) {
	case func(PrivateMessageEvent):
		typedHandler = handler.(func(PrivateMessageEvent)) // 类型断言
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: PrivateMessageEventType,
		})
	case func(GroupMessageEvent):
		typedHandler = handler.(func(GroupMessageEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupMessageEventType,
		})
	case func(TempMessageEvent):
		typedHandler = handler.(func(TempMessageEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: TempMessageEventType,
		})
	case func(NewFriendRequestEvent):
		typedHandler = handler.(func(NewFriendRequestEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: NewFriendRequestEventType,
		})
	case func(NewFriendEvent):
		typedHandler = handler.(func(NewFriendEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: NewFriendEventType,
		})
	case func(FriendRecallEvent):
		typedHandler = handler.(func(FriendRecallEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: FriendRecallEventType,
		})
	case func(FriendRenameEvent):
		typedHandler = handler.(func(FriendRenameEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: FriendRenameEventType,
		})
	case func(FriendPokeEvent):
		typedHandler = handler.(func(FriendPokeEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: FriendPokeEventType,
		})
	case func(GroupMemberPermissionUpdatedEvent):
		typedHandler = handler.(func(GroupMemberPermissionUpdatedEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupMemberPermissionUpdatedEventType,
		})
	case func(GroupNameUpdatedEvent):
		typedHandler = handler.(func(GroupNameUpdatedEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupNameUpdatedEventType,
		})
	case func(GroupMuteEvent):
		typedHandler = handler.(func(GroupMuteEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupMuteEventType,
		})
	case func(GroupRecallEvent):
		typedHandler = handler.(func(GroupRecallEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupRecallEventType,
		})
	case func(GroupMemberJoinRequestEvent):
		typedHandler = handler.(func(GroupMemberJoinRequestEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupMemberJoinRequestEventType,
		})
	case func(GroupMemberIncreaseEvent):
		typedHandler = handler.(func(GroupMemberIncreaseEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupMemberIncreaseEventType,
		})
	case func(GroupMemberDecreaseEvent):
		typedHandler = handler.(func(GroupMemberDecreaseEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupMemberDecreaseEventType,
		})
	case func(GroupDigestEvent):
		typedHandler = handler.(func(GroupDigestEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupDigestEventType,
		})
	case func(GroupReactionEvent):
		typedHandler = handler.(func(GroupReactionEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupReactionEventType,
		})
	case func(GroupMemberSpecialTitleUpdated):
		typedHandler = handler.(func(GroupMemberSpecialTitleUpdated))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupMemberSpecialTitleUpdatedEventType,
		})
	case func(GroupInviteEvent):
		typedHandler = handler.(func(GroupInviteEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: GroupInviteEventType,
		})
	case func(BotConnectedEvent):
		typedHandler = handler.(func(BotConnectedEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: BotConnectedEventType,
		})
	case func(BotDisconnectedEvent):
		typedHandler = handler.(func(BotDisconnectedEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: BotDisconnectedEventType,
		})
	case func(CustomEvent):
		typedHandler = handler.(func(CustomEvent))
		wrapper := TypedWrapper(typedHandler)
		h.Subscriptions = append(h.Subscriptions, Subscription{
			HandlerFunc: wrapper,
			HandlerType: CustomEventType,
		})
	default:
		Warn("传入了不支持的事件类型！")
	}
	return h
}

// HandleMessage 用于向事件处理器添加消息处理函数
func (h *Handler) HandleMessage(handler func(MessageEvent)) *Handler {
	h.MessageMiddlewares = append(h.MessageMiddlewares, func(e CryoEvent) CryoEvent {
		// 尝试对事件进行类型断言
		if msgEvent, ok := e.(CryoMessageEvent); ok {
			handler(msgEvent.GetMessageEvent())
		}
		return e
	})
	return h
}

// Register 将当前的事件处理器注册到事件总线
func (h *Handler) Register() {
	// 将事件处理器中的所有处理函数注册到事件总线
	// 当事件处理器有匹配的事件类型时，只会注册拥有匹配的类型的处理函数
	if len(h.MatchingTypes) == 0 {
		// 如果没有匹配的事件类型，则注册所有的处理函数
		for _, sub := range h.Subscriptions {
			sub.HandlerId = Subscribe(sub.HandlerType, sub.HandlerFunc, h.Tags...)
		}
		// 注册中间件
		AddGlobalMiddleware(h.Middlewares...)
		// 注册消息中间件
		for _, et := range messageEventTypes {
			AddMiddleware(et, h.MessageMiddlewares...)
		}
	} else {
		// 如果有匹配的事件类型，则只注册拥有匹配的类型的处理函数
		for _, matchingType := range h.MatchingTypes { // 遍历所有匹配的事件类型
			// 订阅所有拥有匹配的事件类型的处理函数
			for _, sub := range h.Subscriptions {
				if sub.HandlerType == matchingType {
					sub.HandlerId = Subscribe(sub.HandlerType, sub.HandlerFunc, h.Tags...)
				}
			}
			// 注册中间件
			AddMiddleware(matchingType, h.Middlewares...)
			// 注册消息中间件，只有同时是匹配的事件类型和消息事件类型才会注册
			for _, et := range messageEventTypes {
				if et == matchingType {
					AddMiddleware(et, h.MessageMiddlewares...)
				}
			}

		}
	}
}
