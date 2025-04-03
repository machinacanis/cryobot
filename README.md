# 🧊CryoBot
![Go Badge](https://img.shields.io/badge/Go-1.24%2B-cyan?logo=go)
[![goreportcard](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/machinacanis/cryobot)

🚧锐意开发中...

cryobot 是一个轻量级聊天机器人开发框架，通过嵌入协议实现  [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo)  来实现简单的部署和迁移。

## 特性

- 内嵌协议实现
- 不出意外的话可以单文件部署
- 为多Bot连接设计
- 消息去重 / 负载均衡
- 可启用的Web后台

## 安装

```bash
go get github.com/machinacanis/cryobot
```

## 快速开始

`script/example.go`是一个最小化的聊天机器人示例，展示了如何使用 cryobot 框架登录账号并处理消息。

你可以查看[文档]()以查看完整的框架功能介绍及一个更全面的示例。

```go
// 尚处于开发阶段，API 可能一定100%会有变动
// 仅供参考
package main

import (
	cryo "github.com/machinacanis/cryobot"
	"github.com/machinacanis/cryobot/client"
	"github.com/machinacanis/cryobot/config"
	"github.com/sirupsen/logrus"
)

func main() {
	cryo.Init(config.CryoConfig{
		LogLevel:                     logrus.InfoLevel,
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   true,
	})

	client.ConnectAll()
	select {}
}

```

## 开发进度
- [x] 基本的登录及信息保存功能
- [x] 多Bot连接支持
- [ ] 消息处理

## Thanks！！！

cryobot基于这些开源项目：

- [Lagrange.Core](https://github.com/LagrangeDev/Lagrange.Core) | NTQQ 协议实现
- [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo) | Lagrange.Core 的 Go 语言实现
- [LagrangeGo-Template](https://github.com/ExquisiteCore/LagrangeGo-Template) | LagrangeGo 的模板示例
- [Logrus](https://github.com/sirupsen/logrus) | 优雅的 Go 日志库
- [FreeCache](https://github.com/coocood/freecache) | 高性能的内存缓存库

向这些项目的贡献者们致以最诚挚的感谢！

## 在找兼容Onebot协议的开发框架？

cryobot 是一个通过内嵌的协议实现来连接客户端的开发框架，它是**针对单一平台的使用场景特化**的，如果你想要一个兼容 Onebot 协议的框架，应该看看这些项目：

- [ZeroBot](https://github.com/wdvxdr1123/ZeroBot) | 基于 Onebot 协议的 Golang 机器人开发框架
- [Nonebot2](https://github.com/nonebot/nonebot2) | 跨平台 Python 异步聊天机器人框架