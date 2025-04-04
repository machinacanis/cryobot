package cryobot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
)

type Bot struct {
	initFlag         bool          // 是否初始化完成
	ConnectedClients []*CryoClient // 已连接的bot客户端列表
}

// NewBot 创建一个新的CryoBot实例
func NewBot() *Bot {
	return &Bot{}
}

// Init 初始化cryobot
//
// 可以传入配置项来覆写默认配置，空的配置项会自动使用默认配置
//
// 如果本地配置文件存在，且没有传入配置项，则会自动加载本地配置文件
func (b *Bot) Init(c ...Config) {
	defaultConfig := Config{
		LogLevel:                     logrus.InfoLevel,
		SignServers:                  []string{DefaultSignServer},
		EnableClientAutoSave:         true,
		EnablePrintLogo:              true,
		EnableConnectPrintMiddleware: true,
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   false,
	}
	if len(c) == 0 { // 如果没有传入配置项，则尝试加载本地配置文件
		co, err := ReadCryoConfig()
		if err == nil {
			c = append(c, co)
			Info("已正在加载本地配置文件")
		}
	}
	if len(c) > 0 {
		if c[0].LogLevel != logrus.InfoLevel {
			defaultConfig.LogLevel = c[0].LogLevel
		}
		if c[0].SignServers != nil {
			defaultConfig.SignServers = c[0].SignServers
		}
		if c[0].EnableClientAutoSave {
			defaultConfig.EnableClientAutoSave = c[0].EnableClientAutoSave
		}
		if c[0].EnablePrintLogo {
			defaultConfig.EnablePrintLogo = c[0].EnablePrintLogo
		}
		if c[0].EnableConnectPrintMiddleware {
			defaultConfig.EnableConnectPrintMiddleware = c[0].EnableConnectPrintMiddleware
		}
		if c[0].EnableMessagePrintMiddleware {
			defaultConfig.EnableMessagePrintMiddleware = c[0].EnableMessagePrintMiddleware
		}
		if c[0].EnableEventDebugMiddleware {
			defaultConfig.EnableEventDebugMiddleware = c[0].EnableEventDebugMiddleware
		}
	}
	conf = defaultConfig // 初始化配置

	// 设置日志等级
	InitTextLogger(conf.LogLevel)
	// 初始化事件总线
	fmt.Print(logo)
	Infof("%s[Cryo] 🧊cryobot 正在初始化...", lavender)
	Bus = NewEventBus() // 初始化事件总线
	// 设置连接打印中间件
	setConnectPrintMiddleware()
	// 设置消息打印中间件
	setMessagePrintMiddleware()
	// 设置事件调试中间件
	setEventDebugMiddleware()

	b.initFlag = true
}

// Start 启动cryobot
func (b *Bot) Start() {
	if !b.initFlag {
		// 没有进行初始化
		log.Fatal("cryobot 没有进行初始化，请先调用 Init() 函数进行初始化！")
	}
	select {} // 阻塞主线程，运行事件循环
}

// AutoConnect 自动连接
func (b *Bot) AutoConnect() {
	if !b.initFlag {
		// 没有进行初始化
		log.Fatal("cryobot 没有进行初始化，请先调用 Init() 函数进行初始化！")
	}
	// 首先检测是否已经连接
	if len(b.ConnectedClients) > 0 {
		// 跳过自动连接
		return
	}
	// 尝试连接所有已保存的bot客户端
	b.ConnectAllSavedClient()
	// 如果没有连接成功，则尝试连接新的bot客户端
	retriedCount := 0
	for len(b.ConnectedClients) == 0 && retriedCount < 3 {
		b.ConnectNewClient()
		retriedCount++
	}
	if len(b.ConnectedClients) == 0 {
		log.Fatal("达到最大重试次数，cryobot 无法连接到bot客户端，请检查网络或配置文件")
	}
}

// ConnectSavedClient 尝试查询并连接到指定的bot客户端
func (b *Bot) ConnectSavedClient(info CryoClientInfo) bool {
	c := NewCryoClient()
	c.Init()
	if !c.Rebuild(info) {
		return false
	}
	Infof("%s[Cryo] 正在连接 %s：%s (%d)", lavender, c.Nickname, c.Id, c.Uin)
	if !c.SignatureLogin() {
		return false
	}
	b.ConnectedClients = append(b.ConnectedClients, c)
	return true
}

// ConnectNewClient 尝试连接一个新的bot客户端
func (b *Bot) ConnectNewClient() bool {
	c := NewCryoClient()
	c.Init()
	Infof("%s[Cryo] 正在连接 %s：%s (%d)", lavender, c.Nickname, c.Id, c.Uin)
	if !c.QRCodeLogin() {
		return false
	}
	b.ConnectedClients = append(b.ConnectedClients, c)
	return true
}

// ConnectAllSavedClient 尝试连接所有已保存的bot客户端
func (b *Bot) ConnectAllSavedClient() {
	// 读取历史连接的客户端
	clientInfos, err := ReadClientInfos()
	if err != nil {
		Error("读取Bot信息时出现错误：", err)
		return
	}
	if len(clientInfos) == 0 {
		Info("没有找到Bot信息")
		return
	}
	for _, info := range clientInfos {
		if !b.ConnectSavedClient(info) {
			Error("通过历史记录连接Bot客户端失败")
			Error("已自动清除失效的客户端信息，请重新登录")
		}
	}
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
	messageEventTypes := []CryoEventType{
		PrivateMessageEventType,
		GroupMessageEventType,
		TempMessageEventType,
	}
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
