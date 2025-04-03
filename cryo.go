package cryo

import (
	"github.com/machinacanis/cryobot/client"
	"github.com/machinacanis/cryobot/config"
	"github.com/machinacanis/cryobot/event"
	"github.com/machinacanis/cryobot/log"
	"github.com/sirupsen/logrus"
)

// Init 初始化cryobot
//
// 可以传入配置项来覆写默认配置，空的配置项会自动使用默认配置
//
// 如果本地配置文件存在，且没有传入配置项，则会自动加载本地配置文件
func Init(c ...config.CryoConfig) {
	defaultConfig := config.CryoConfig{
		LogLevel:                     logrus.InfoLevel,
		SignServers:                  []string{config.DefaultSignServer},
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
	// 设置消息打印中间件
	setMessagePrintMiddleware()
	// 设置事件调试中间件
	setEventDebugMiddleware()
}

// setMessagePrintMiddleware 内置的消息打印中间件
func setMessagePrintMiddleware() {
	if config.Conf.EnableMessagePrintMiddleware {
		client.Bus.UseMiddleware(event.PrivateMessageEventType, func(e event.CryoEvent) event.CryoEvent {
			if typedEvent, ok := e.(event.PrivateMessageEvent); ok {
				log.Infof("[私聊][%d] From %s(%d)   %s", typedEvent.SelfUin, typedEvent.SenderNickname, typedEvent.SenderUin, typedEvent.MessageElements.ToString())
			}
			return e
		})
		client.Bus.UseMiddleware(event.GroupMessageEventType, func(e event.CryoEvent) event.CryoEvent {
			if typedEvent, ok := e.(event.GroupMessageEvent); ok {
				log.Infof("[群 %s(%d)][%d] From %s(%d)   %s", typedEvent.GroupName, typedEvent.GroupUin, typedEvent.SelfUin, typedEvent.SenderNickname, typedEvent.SenderUin, typedEvent.MessageElements.ToString())
			}
			return e
		})
	}
}

// setEventDebugMiddleware 内置的事件调试中间件
func setEventDebugMiddleware() {
	if config.Conf.EnableEventDebugMiddleware {
		client.Bus.UseGlobalMiddleware(func(e event.CryoEvent) event.CryoEvent {
			log.Debug(e.ToJsonString())
			return e
		})
	}
}

func Start() {}

func Stop() {}

func Restart() {}
