package cryo

import (
	"github.com/machinacanis/cryobot/client"
	"github.com/machinacanis/cryobot/config"
	"github.com/machinacanis/cryobot/event"
	"github.com/machinacanis/cryobot/log"
	"github.com/sirupsen/logrus"
)

var initFlag = false

// Init 初始化cryobot
//
// 可以传入配置项来覆写默认配置，空的配置项会自动使用默认配置
//
// 如果本地配置文件存在，且没有传入配置项，则会自动加载本地配置文件
func Init(c ...config.CryoConfig) {
	defaultConfig := config.CryoConfig{
		LogLevel:                     logrus.InfoLevel,
		SignServers:                  []string{config.DefaultSignServer},
		EnableConnectPrintMiddleware: true,
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   false,
	}
	if len(c) == 0 { // 如果没有传入配置项，则尝试加载本地配置文件
		co, err := config.ReadCryoConfig()
		if err == nil {
			c = append(c, co)
			log.Info("已正在加载本地配置文件")
		}
	}
	if len(c) > 0 {
		if c[0].LogLevel != logrus.InfoLevel {
			defaultConfig.LogLevel = c[0].LogLevel
		}
		if c[0].SignServers != nil {
			defaultConfig.SignServers = c[0].SignServers
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
	config.Conf = defaultConfig // 初始化配置

	// 设置日志等级
	log.InitTextLogger(config.Conf.LogLevel)
	// 初始化事件总线
	log.Info("🧊cryobot 正在初始化...")
	client.Bus = client.NewEventBus() // 初始化事件总线
	// 设置连接打印中间件
	setConnectPrintMiddleware()
	// 设置消息打印中间件
	setMessagePrintMiddleware()
	// 设置事件调试中间件
	setEventDebugMiddleware()

	initFlag = true
}

// setConnectPrintMiddleware 内置的连接打印中间件
func setConnectPrintMiddleware() {
	if config.Conf.EnableConnectPrintMiddleware {
		client.Bus.UseMiddleware(event.BotConnectedEventType, func(e event.CryoEvent) event.CryoEvent {
			if typedEvent, ok := e.(event.BotConnectedEvent); ok {
				log.Infof("%s：%s (%d) 已成功连接", typedEvent.BotNickname, typedEvent.BotId, typedEvent.BotUin)
			}
			return e
		})
		client.Bus.UseMiddleware(event.BotDisconnectedEventType, func(e event.CryoEvent) event.CryoEvent {
			if typedEvent, ok := e.(event.BotDisconnectedEvent); ok {
				log.Infof("%s：%s (%d) 已成功连接", typedEvent.BotNickname, typedEvent.BotId, typedEvent.BotUin)
			}
			return e
		})
	}
}

// setMessagePrintMiddleware 内置的消息打印中间件
func setMessagePrintMiddleware() {
	if config.Conf.EnableMessagePrintMiddleware {
		client.Bus.UseMiddleware(event.PrivateMessageEventType, func(e event.CryoEvent) event.CryoEvent {
			if typedEvent, ok := e.(event.PrivateMessageEvent); ok {
				log.Infof("[%s][私聊] From %s(%d) - %s", typedEvent.BotNickname, typedEvent.SenderNickname, typedEvent.SenderUin, typedEvent.MessageElements.ToString())
			}
			return e
		})
		client.Bus.UseMiddleware(event.GroupMessageEventType, func(e event.CryoEvent) event.CryoEvent {
			if typedEvent, ok := e.(event.GroupMessageEvent); ok {
				log.Infof("[%s][群 %s(%d)] From %s(%d) - %s", typedEvent.BotNickname, typedEvent.GroupName, typedEvent.GroupUin, typedEvent.SenderNickname, typedEvent.SenderUin, typedEvent.MessageElements.ToString())
			}
			return e
		})
	}
}

// setEventDebugMiddleware 内置的事件调试中间件
func setEventDebugMiddleware() {
	if config.Conf.EnableEventDebugMiddleware {
		client.Bus.UseGlobalMiddleware(func(e event.CryoEvent) event.CryoEvent {
			log.Debug(e.ToJsonString()) // 输出json
			return e
		})
	}
}

// Start 启动cryobot
func Start(c ...config.CryoConfig) {
	if !initFlag {
		// 没有进行初始化
		log.Fatal("cryobot 没有进行初始化，请先调用 Init() 函数进行初始化！")
	}
	select {} // 阻塞主线程，运行事件循环
}

// AutoConnect 自动连接bot客户端，如果已经连接则跳过，未连接则尝试从连接历史中建立连接，如果全部失败则尝试创建新的bot客户端
func AutoConnect() {
	if !initFlag {
		// 没有进行初始化
		log.Fatal("cryobot 没有进行初始化，请先调用 Init() 函数进行初始化！")
	}
	// 首先检测是否已经连接
	if len(client.ConnectedClients) > 0 {
		// 跳过自动连接
		return
	}
	// 尝试连接所有已保存的bot客户端
	client.ConnectAll()
	// 如果没有连接成功，则尝试连接新的bot客户端
	retriedCount := 0
	for len(client.ConnectedClients) == 0 && retriedCount < 3 {
		client.ConnectNew()
		retriedCount++
	}
	if len(client.ConnectedClients) == 0 {
		log.Fatal("达到最大重试次数，cryobot 无法连接到bot客户端，请检查网络或配置文件")
	}
}
