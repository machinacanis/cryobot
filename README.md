# CryoBot

锐意开发中...

cryobot 是一个轻量级聊天机器人开发框架，直接嵌入了协议实现  [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo)  ，使部署和迁移变得无比简单。

## 特性

- 内嵌  [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo)  协议实现
- 极简的单文件部署
- 良好的多Bot连接支持
- 消息去重 / 负载均衡
- 可启用的Web后台

## 安装

```bash
go get github.com/machinacanis/cryobot
```

## 快速开始

`main.go`是一个最小化的聊天机器人示例，展示了如何使用 cryobot 框架登录账号并处理消息。

你可以查看[文档]()以查看完整的框架功能介绍及一个更全面的示例。

```go
暂未完成
```

## 开发进度
- [x] 基本的登录及信息保存功能
- [x] 多Bot连接支持
- [ ] 消息处理

## 鸣谢

- [Lagrange.Core](https://github.com/LagrangeDev/Lagrange.Core)
- [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo)
- [LagrangeGo-Template](https://github.com/ExquisiteCore/LagrangeGo-Template)