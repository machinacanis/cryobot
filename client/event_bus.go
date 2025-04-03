/*

event_bus.go

cryobot的事件总线，基于简单的发布-订阅架构，支持中间件和异步处理

*/

package client

import (
	"github.com/machinacanis/cryobot/event"
	"sync"
)

var Bus *EventBus

// EventHandler 是一个函数类型，用于处理事件
type EventHandler func(event event.CryoEvent)

// Middleware 是一个函数类型，用于定义事件处理过程中的中间件函数
type Middleware func(event event.CryoEvent) event.CryoEvent

// EventBus 是一个事件总线，用于管理事件的订阅和发布
type EventBus struct {
	mutex       sync.RWMutex
	subscribers map[event.CryoEventType][]EventHandler
	middleware  map[event.CryoEventType][]Middleware
}

// NewEventBus 创建一个新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[event.CryoEventType][]EventHandler),
		middleware:  make(map[event.CryoEventType][]Middleware),
	}
}

// UseMiddleware 注册一个中间件，用于处理特定事件类型的事件
func (bus *EventBus) UseMiddleware(eventType event.CryoEventType, middleware Middleware) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	if _, exists := bus.middleware[eventType]; !exists {
		bus.middleware[eventType] = []Middleware{}
	}
	bus.middleware[eventType] = append(bus.middleware[eventType], middleware)
}

// UseGlobalMiddleware 注册一个全局中间件，用于处理所有事件类型的事件
func (bus *EventBus) UseGlobalMiddleware(middleware Middleware) {
	for _, eventType := range AllEventTypes() {
		bus.UseMiddleware(eventType, middleware)
	}
}

// ApplyMiddleware 应用中间件
func (bus *EventBus) applyMiddleware(event event.CryoEvent) event.CryoEvent {
	eventType := event.Type()

	bus.mutex.RLock()
	middles, exists := bus.middleware[eventType]
	bus.mutex.RUnlock()

	if !exists || len(middles) == 0 {
		return event
	}

	result := event
	for _, middleware := range middles {
		if result == nil {
			return nil // Event was blocked by middleware
		}
		result = middleware(result)
	}

	return result
}

// Subscribe 注册一个处理器以接收特定事件类型的事件
func (bus *EventBus) Subscribe(eventType event.CryoEventType, handler EventHandler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	if _, exists := bus.subscribers[eventType]; !exists {
		bus.subscribers[eventType] = []EventHandler{}
	}
	bus.subscribers[eventType] = append(bus.subscribers[eventType], handler)
}

// Unsubscribe 移除一个处理器，使其不再接收特定事件类型的事件
func (bus *EventBus) Unsubscribe(eventType event.CryoEventType, handler EventHandler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	handlers, exists := bus.subscribers[eventType]
	if !exists {
		return
	}

	// 找到并移除处理器
	for i, h := range handlers {
		if &h == &handler { // 比较函数地址
			bus.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// UnsubscribeAll 移除所有处理器，使其不再接收特定事件类型的事件，危险操作，慎用
func (bus *EventBus) UnsubscribeAll(eventType event.CryoEventType) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	// 移除所有处理器
	if _, exists := bus.subscribers[eventType]; exists {
		delete(bus.subscribers, eventType)
	}
}

// Publish 发布一个事件，所有订阅该事件类型的处理器将被调用
func (bus *EventBus) Publish(event event.CryoEvent) {
	// 应用中间件
	processedEvent := bus.applyMiddleware(event)
	if processedEvent == nil {
		return // 事件被中间件阻止
	}

	eventType := processedEvent.Type()

	bus.mutex.RLock()
	handlers, exists := bus.subscribers[eventType]
	bus.mutex.RUnlock()

	if !exists {
		return
	}

	// 触发所有处理器
	for _, handler := range handlers {
		handler(processedEvent)
	}
}

func (bus *EventBus) AsyncPublish(event event.CryoEvent) {
	// 应用中间件
	processedEvent := bus.applyMiddleware(event)
	if processedEvent == nil {
		return // 事件被中间件阻止
	}

	eventType := processedEvent.Type()

	bus.mutex.RLock()
	handlers, exists := bus.subscribers[eventType]
	bus.mutex.RUnlock()

	if !exists || len(handlers) == 0 {
		return
	}

	// 如果只有一个处理器，不并发
	if len(handlers) == 1 {
		handlers[0](processedEvent)
		return
	}
	// 创建一个等待组来等待所有处理器运行
	var wg sync.WaitGroup
	wg.Add(len(handlers))
	// 每个处理器在自己的 goroutine 中启动
	for _, handler := range handlers {
		h := handler
		e := processedEvent
		go func() {
			defer wg.Done()
			h(e)
		}()
	}
	// 射 后 不 理
	// wg.Wait()
}

// SubscribeAll 对所有事件类型注册一个处理器，没事别乱用
func (bus *EventBus) SubscribeAll(handler EventHandler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	// Get all available event types
	for _, eventType := range AllEventTypes() {
		if _, exists := bus.subscribers[eventType]; !exists {
			bus.subscribers[eventType] = []EventHandler{}
		}
		bus.subscribers[eventType] = append(bus.subscribers[eventType], handler)
	}
}

// AllEventTypes 返回所有可用的事件类型
func AllEventTypes() []event.CryoEventType {
	return []event.CryoEventType{
		event.PrivateMessageEventType,
		event.GroupMessageEventType,
		event.TempMessageEventType,
		event.NewFriendRequestEventType,
		event.NewFriendEventType,
		event.FriendRecallEventType,
		event.FriendRenameEventType,
		event.FriendPokeEventType,
		event.GroupMemberPermissionUpdatedEventType,
		event.GroupNameUpdatedEventType,
		event.GroupMuteEventType,
		event.GroupRecallEventType,
		event.GroupMemberJoinRequestEventType,
		event.GroupMemberIncreaseEventType,
		event.GroupMemberDecreaseEventType,
		event.GroupDigestEventType,
		event.GroupReactionEventType,
		event.GroupMemberSpecialTitleUpdatedEventType,
		event.GroupInviteEventType,
		event.BotConnectedEventType,
		event.BotDisconnectedEventType,
		event.CustomEventType,
	}
}

// SubscribeSpecific 订阅特定类型的事件，并将处理器转换为特定类型
func SubscribeSpecific[T event.CryoEvent](eventType event.CryoEventType, handler func(T)) {
	wrapperHandler := func(event event.CryoEvent) {
		if typedEvent, ok := event.(T); ok {
			handler(typedEvent)
		}
	}
	Bus.Subscribe(eventType, wrapperHandler)
}

// Subscribe 订阅特定类型的事件
func Subscribe(eventType event.CryoEventType, handler func(event event.CryoEvent)) {
	Bus.Subscribe(eventType, handler)
}

// Unsubscribe 取消订阅特定类型的事件
func Unsubscribe(eventType event.CryoEventType, handler func(event event.CryoEvent)) {
	Bus.Unsubscribe(eventType, handler)
}

// SubscribeAll 订阅所有事件类型，尽量不要滥用，相当于啥事件都要走一遍这个handler，不太好说整多了会不会产生性能问题
func SubscribeAll(handler func(event event.CryoEvent)) {
	Bus.SubscribeAll(handler)
}

// UnsubscribeAll 移除所有指定类型的处理器，使其不再接收事件，高危操作，慎用
func UnsubscribeAll(eventType event.CryoEventType) {
	Bus.UnsubscribeAll(eventType)
}

// Publish 发布一个事件
func Publish(event event.CryoEvent) {
	Bus.Publish(event)
}

// AsyncPublish 异步发布一个事件
func AsyncPublish(event event.CryoEvent) {
	Bus.AsyncPublish(event)
}

// UseMiddleware 注册一个中间件，用于处理特定事件类型的事件
func UseMiddleware(eventType event.CryoEventType, middleware Middleware) {
	Bus.UseMiddleware(eventType, middleware)
}

// UseGlobalMiddleware 注册一个全局中间件，用于处理所有事件类型的事件
func UseGlobalMiddleware(middleware Middleware) {
	Bus.UseGlobalMiddleware(middleware)
}
