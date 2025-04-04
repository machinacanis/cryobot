package cryobot

import "sync"

var Bus *CryoEventBus

// CryoEventHandler 用于实现事件处理器的接口
type CryoEventHandler interface {
	GetType() CryoEventType
	GetId() string
	GetTags() []string // 新增：获取处理器标签
	Handle(event CryoEvent)
}

// EventHandler 是一个泛型事件处理器，用于处理特定类型的事件
type EventHandler[T CryoEvent] struct {
	handlerType CryoEventType
	handlerId   string
	handler     func(event T)
	tags        []string
}

// GetType 返回事件处理器支持的事件类型
func (h EventHandler[T]) GetType() CryoEventType {
	return h.handlerType
}

// GetId 返回事件处理器的唯一标识符
func (h EventHandler[T]) GetId() string {
	return h.handlerId
}

// Handle 处理事件
func (h EventHandler[T]) Handle(event CryoEvent) {
	// 使用类型断言来确保事件类型匹配
	if typedEvent, ok := event.(T); ok {
		h.handler(typedEvent)
	}
}

// GetTags 返回事件处理器的标签
func (h EventHandler[T]) GetTags() []string {
	return h.tags
}

// Middleware 是一个函数类型，用于定义事件处理过程中的中间件函数
type Middleware func(event CryoEvent) CryoEvent

// CryoEventBus 是一个事件总线，用于管理事件的订阅和发布
type CryoEventBus struct {
	subscriberMutex sync.RWMutex
	middlewareMutex sync.RWMutex
	subscriber      map[CryoEventType][]CryoEventHandler
	middleware      map[CryoEventType][]Middleware
}

// NewEventBus 创建一个新的事件总线
func NewEventBus() *CryoEventBus {
	return &CryoEventBus{
		subscriber: make(map[CryoEventType][]CryoEventHandler),
		middleware: make(map[CryoEventType][]Middleware),
	}
}

// applyMiddleware 应用中间件
func (bus *CryoEventBus) applyMiddleware(event CryoEvent) CryoEvent {
	eventType := event.Type()

	bus.middlewareMutex.RLock()
	middlewareSlice, exists := bus.middleware[eventType]
	// 创建一个中间件切片的副本，以避免在应用中间件时发生并发修改
	var middlewareCopy []Middleware
	if exists {
		middlewareCopy = make([]Middleware, len(middlewareSlice))
		copy(middlewareCopy, middlewareSlice)
	}
	bus.middlewareMutex.RUnlock()

	if !exists {
		return event
	}

	// 优化了一下，在不持有锁的时候应用中间件
	currentEvent := event
	for _, middleware := range middlewareCopy {
		currentEvent = middleware(currentEvent)
		if currentEvent == nil {
			return nil
		}
	}
	return currentEvent
}

// Subscribe 注册一个事件处理器，用于处理特定类型的事件
func Subscribe[T CryoEvent](eventType CryoEventType, handler func(event T), tag ...string) string {
	Bus.subscriberMutex.Lock()
	defer Bus.subscriberMutex.Unlock()

	if _, exists := Bus.subscriber[eventType]; !exists {
		Bus.subscriber[eventType] = []CryoEventHandler{}
	}

	// 如果提供了标签，使用标签；否则使用空字符串
	var handlerTag []string
	if len(tag) > 0 {
		handlerTag = tag
	}

	// 生成唯一标识符
	handlerId := NewUUID()

	eventHandler := EventHandler[T]{
		handlerType: eventType,
		handlerId:   handlerId,
		tags:        handlerTag,
		handler:     handler,
	}

	Bus.subscriber[eventType] = append(Bus.subscriber[eventType], eventHandler)

	// 返回handlerId，以便用户可以选择使用id或tag来解除订阅
	return handlerId
}

// Publish 发布事件
func Publish(event CryoEvent) {
	Bus.subscriberMutex.RLock()
	defer Bus.subscriberMutex.RUnlock()

	eventType := event.Type()

	// 应用中间件
	processedEvent := Bus.applyMiddleware(event)
	if processedEvent == nil {
		return // 事件被中间件截断
	}

	// 获取订阅内容
	handlers, exists := Bus.subscriber[eventType]
	if !exists {
		return
	}

	// 依次调用处理器
	for _, handler := range handlers {
		handler.Handle(processedEvent)
	}
}

// PublishAsync 异步发布事件
func PublishAsync(event CryoEvent) {
	eventType := event.Type()

	// 应用中间件
	processedEvent := Bus.applyMiddleware(event)
	if processedEvent == nil {
		return
	}

	// 复制处理器列表
	Bus.subscriberMutex.RLock()
	handlers, exists := Bus.subscriber[eventType]
	if !exists {
		Bus.subscriberMutex.RUnlock()
		return
	}

	// 避免在锁内进行处理器调用，创建一个处理器列表的副本
	handlersCopy := make([]CryoEventHandler, len(handlers))
	copy(handlersCopy, handlers)
	Bus.subscriberMutex.RUnlock()

	// 使用 goroutine 异步调用处理器
	go func() {
		var wg sync.WaitGroup
		for _, handler := range handlersCopy {
			wg.Add(1)
			go func(h CryoEventHandler, e CryoEvent) {
				defer wg.Done()
				h.Handle(e)
			}(handler, processedEvent)
		}
		wg.Wait()
	}()
}

// AddMiddleware 为特定事件类型添加中间件
func AddMiddleware(eventType CryoEventType, middleware ...Middleware) {
	Bus.middlewareMutex.Lock()
	defer Bus.middlewareMutex.Unlock()

	if _, exists := Bus.middleware[eventType]; !exists {
		Bus.middleware[eventType] = []Middleware{}
	}
	Bus.middleware[eventType] = append(Bus.middleware[eventType], middleware...)
}

// AddGlobalMiddleware 为所有事件类型添加中间件
func AddGlobalMiddleware(middleware ...Middleware) {
	Bus.middlewareMutex.Lock()
	defer Bus.middlewareMutex.Unlock()

	for eventType := range Bus.subscriber {
		if _, exists := Bus.middleware[eventType]; !exists {
			Bus.middleware[eventType] = []Middleware{}
		}
		Bus.middleware[eventType] = append(Bus.middleware[eventType], middleware...)
	}
}

// UnsubscribeById 取消订阅事件处理器
func UnsubscribeById(handlerId string) {
	Bus.subscriberMutex.Lock()
	defer Bus.subscriberMutex.Unlock()

	for eventType, handlers := range Bus.subscriber {
		for i, handler := range handlers {
			if handler.GetId() == handlerId {
				// 删除处理器
				Bus.subscriber[eventType] = append(handlers[:i], handlers[i+1:]...)
				return
			}
		}
	}
}

// UnsubscribeByTag 取消订阅事件处理器
func UnsubscribeByTag(tag ...string) {
	if len(tag) == 0 {
		return
	}

	Bus.subscriberMutex.Lock()
	defer Bus.subscriberMutex.Unlock()

	for eventType, handlers := range Bus.subscriber {
		// 创建新的切片来存储不包含指定标签的处理器
		newHandlers := make([]CryoEventHandler, 0, len(handlers))
		for _, handler := range handlers {
			if !containsAllTags(handler.GetTags(), tag) {
				newHandlers = append(newHandlers, handler)
			}
		}

		// 只在处理器数量发生变化时更新
		if len(newHandlers) != len(handlers) {
			Bus.subscriber[eventType] = newHandlers
		}
	}
}

// containsAllTags 检查处理器是否包含所有传入的标签
func containsAllTags(handlerTags, tags []string) bool {
	for _, tag := range tags {
		found := false
		for _, handlerTag := range handlerTags {
			if handlerTag == tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
