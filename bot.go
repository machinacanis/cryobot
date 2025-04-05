package cryobot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
)

type Bot struct {
	initFlag         bool                   // 是否初始化完成
	ConnectedClients map[string]*CryoClient // 已连接的Bot客户端集合
}

// NewBot 创建一个新的CryoBot实例
//
// 非常不推荐同时创建多个Bot实例，因为事件总线是全局的，如果你实在要这么做那么请确保你知道自己在做什么
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
	// 初始化连接的客户端集合
	b.ConnectedClients = make(map[string]*CryoClient)
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
	b.ConnectedClients[c.Id] = c
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
	b.ConnectedClients[c.Id] = c
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

// GetClientById 获取指定ID的bot客户端
func (b *Bot) GetClientById(id string) *CryoClient {
	if client, ok := b.ConnectedClients[id]; ok {
		return client
	}
	return nil
}

// GetClientByUin 获取指定Uin的bot客户端
func (b *Bot) GetClientByUin(uin int) *CryoClient {
	for _, client := range b.ConnectedClients {
		if client.Uin == uin {
			return client
		}
	}
	return nil
}

// GetClientByUid 获取指定Uid的bot客户端
func (b *Bot) GetClientByUid(uid string) *CryoClient {
	for _, client := range b.ConnectedClients {
		if client.Uid == uid {
			return client
		}
	}
	return nil
}

// GetClient 获取指定事件对应的bot客户端
func (b *Bot) GetClient(event CryoEvent) *CryoClient {
	return b.GetClientById(event.GetBaseEvent().BotId)
}

func (b *Bot) Send(event CryoMessageEvent, args ...interface{}) (ok bool, messageId uint32) {
	// 根据事件获取对应的bot客户端
	return b.GetClient(event).Send(event, args...)
}

func (b *Bot) Reply(event CryoMessageEvent, args ...interface{}) (ok bool, messageId uint32) {
	// 根据事件获取对应的bot客户端
	return b.GetClient(event).Reply(event, args...)
}
