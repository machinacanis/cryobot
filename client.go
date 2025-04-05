package cryobot

import (
	"encoding/base64"
	"fmt"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/auth"
	"os"
	"time"
)

// CryoClient cryobot的Bot客户端封装
type CryoClient struct {
	Id        string
	Client    *client.QQClient
	Platform  string
	Version   string
	DeviceNum int
	Uin       int
	Uid       string
	Nickname  string

	initFlag bool // 是否初始化完成
}

// NewCryoClient 创建一个新的CryoClient实例
func NewCryoClient() *CryoClient {
	return &CryoClient{}
}

// Init 初始化一个新的CryoClient客户端
func (c *CryoClient) Init() {
	c.Id = NewUUID() // 给Bot客户端分配一个唯一的UUID

	// 默认平台和版本
	if c.Platform == "" {
		c.Platform = "linux"
	}
	if c.Version == "" {
		c.Version = "3.2.15-30366"
	}

	appInfo := auth.AppList[c.Platform][c.Version]
	c.Client = client.NewClient(0, "")
	c.Client.SetLogger(pLogger) // 替换日志记录器，详见client/protocol_logger.go以及log/logger.go
	c.Client.UseVersion(appInfo)
	c.Client.AddSignServer(conf.SignServers...)
	c.DeviceNum = RandomDeviceNumber()
	c.Client.UseDevice(auth.NewDeviceInfo(c.DeviceNum))
	c.Nickname = newNickname() // 生成一个默认的编号昵称

	c.initFlag = true
}

// Rebuild 重新构建CryoClient实例
func (c *CryoClient) Rebuild(clientInfo CryoClientInfo) bool {
	if !c.initFlag {
		Error("cryobot客户端没有完成初始化，请先调用Init()方法")
		return false
	}
	var sig string
	c.Id = clientInfo.Id
	c.Platform = clientInfo.Platform
	c.Version = clientInfo.Version
	c.DeviceNum = clientInfo.DeviceNum
	c.Uin = clientInfo.Uin
	c.Uid = clientInfo.Uid
	sig = clientInfo.Signature
	c.Client.UseDevice(auth.NewDeviceInfo(c.DeviceNum))
	c.Client.UseVersion(auth.AppList[c.Platform][c.Version])
	c.UseSignature(sig) // 使用指定的签名信息
	return true
}

// Save 将当前客户端的信息保存到文件中
func (c *CryoClient) Save() error {
	clientInfo := CryoClientInfo{
		Id:        c.Id,
		Signature: c.GetSignature(),
		Platform:  c.Platform,
		Version:   c.Version,
		DeviceNum: c.DeviceNum,
		Uin:       c.Uin,
		Uid:       c.Uid,
	}
	err := SaveClientInfo(clientInfo)
	return err
}

// GetSignature 获取当前客户端的签名信息
func (c *CryoClient) GetSignature() string {
	data, err := c.Client.Sig().Marshal()
	if err != nil {
		Error("序列化签名时出现错误：", err)
		return ""
	}
	// 将二进制的签名直接编码到字符串
	sig := base64.StdEncoding.EncodeToString(data)
	return sig
}

// UseSignature 使用指定的签名信息
func (c *CryoClient) UseSignature(sig string) {
	// 将字符串解码为二进制
	data, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		Error("解码签名时出现错误：", err)
		return
	}
	// 反序列化签名
	sigInfo, err := auth.UnmarshalSigInfo(data, true)
	if err != nil {
		Error("反序列化签名时出现错误：", err)
		return
	}
	c.Client.UseSig(sigInfo)
}

func (c *CryoClient) AfterLogin() {
	// 登录成功后，保存签名
	c.Uin = int(c.Client.Sig().Uin)
	c.Uid = c.Client.Sig().UID
	SendBotConnectedEvent(c)       // 发送登录成功事件
	if conf.EnableClientAutoSave { // 如果启用了自动保存
		err := c.Save()
		if err != nil {
			Error("保存登录信息时出现错误：", err)
		} // 保存登录信息
	}

	// 订阅事件
	EventBind(c)
}

// GetQRCode 获取二维码信息
func (c *CryoClient) GetQRCode() ([]byte, string, error) {
	code, res, err := c.Client.FetchQRCodeDefault()
	// 这里获取到两个参数，第一个是字节形式的二维码图片，第二个是二维码指向的链接
	return code, res, err
}

// SaveQRCode 保存二维码图片
func (c *CryoClient) SaveQRCode(code []byte) bool {
	qrcodePath := fmt.Sprintf("QRCode_%s.png", c.Id)
	err := os.WriteFile(qrcodePath, code, 0644)
	if err != nil {
		Error("写入二维码图片时出现错误：", err)
		return false
	}
	Infof("登录二维码已保存到 %s", qrcodePath)
	return true
}

// PrintQRCode 打印二维码
func (c *CryoClient) PrintQRCode(url string) {
	// 打印二维码的链接
	fmt.Println(*GetQRCodeString(url)) // 注意使用了指针
}

// SignatureLogin 使用签名快速登录
func (c *CryoClient) SignatureLogin() (ok bool) {
	sig := c.Client.Sig()
	if sig != nil {
		err := c.Client.FastLogin()
		if err == nil {
			// 通过保存的签名快速登录成功
			c.AfterLogin()
			return true
		}
	}
	return false
}

// QRCodeLogin 使用二维码登录
func (c *CryoClient) QRCodeLogin() bool {
	Info("正在使用二维码登录...")
	code, url, err := c.GetQRCode()
	if err != nil {
		Error("获取二维码时出现错误：", err)
		return false
	}
	// 保存二维码图片
	c.SaveQRCode(code)
	// 向终端输出二维码
	c.PrintQRCode(url)
	if !c.watingForLoginResult() { // 等待扫码登录
		Warn("扫码登录失败！")
	}
	c.AfterLogin()
	return true
}

// watingForLoginResult 等待扫码登录结果
func (c *CryoClient) watingForLoginResult() bool {
	//轮询登录状态
	for {
		retCode, err := c.Client.GetQRCodeResult()
		if err != nil {
			Error("获取二维码登录结果时出现错误：", err)
			return false
		}
		// 等待扫码
		if retCode.Waitable() {
			time.Sleep(1 * time.Second)
			continue
		}
		if !retCode.Success() {
			return false
		}
		break
	}
	_, err := c.Client.QRCodeLogin()
	if err != nil {
		Error("二维码登录时出现错误：", err)
		return false
	}
	return true
}

// SendPrivateMessage 发送私聊消息
func (c *CryoClient) SendPrivateMessage(userUin uint32, msg *CryoMessage) (ok bool, messageId uint32) {
	// 发送私聊消息
	message, err := c.Client.SendPrivateMessage(userUin, msg.ToIMessageElements())
	if err != nil {
		Errorf("向用户 %d 发送消息时出现错误：%v", userUin, err)
		return false, 0
	}
	return true, message.ID
}

// SendGroupMessage 发送群消息
func (c *CryoClient) SendGroupMessage(groupUin uint32, msg *CryoMessage) (ok bool, messageId uint32) {
	// 发送群消息
	message, err := c.Client.SendGroupMessage(groupUin, msg.ToIMessageElements())
	if err != nil {
		Errorf("向群 %d 发送消息时出现错误：%v", groupUin, err)
		return false, 0
	}
	return true, message.ID
}

// SendTempMessage 发送临时消息
func (c *CryoClient) SendTempMessage(groupUin, userUin uint32, msg *CryoMessage) (ok bool, messageId uint32) {
	// 发送临时消息
	message, err := c.Client.SendTempMessage(groupUin, userUin, msg.ToIMessageElements())
	if err != nil {
		Errorf("向与用户 %d 的临时会话发送消息时出现错误：%v", groupUin, err)
		return false, 0
	}
	return true, message.ID
}

func (c *CryoClient) Send(event CryoMessageEvent, args ...interface{}) (ok bool, messageId uint32) {
	// 处理消息内容
	m := ProcessMessageContent(args...)
	// 根据传入的事件来发送消息
	switch event.Type() {
	case PrivateMessageEventType:
		// 断言为私聊消息事件
		if msg, ok := event.(PrivateMessageEvent); ok {
			return c.SendPrivateMessage(msg.SenderUin, m)
		}
	case GroupMessageEventType:
		// 断言为群消息事件
		if msg, ok := event.(GroupMessageEvent); ok {
			return c.SendGroupMessage(msg.GroupUin, m)
		}
	case TempMessageEventType:
		// 断言为临时消息事件
		if msg, ok := event.(TempMessageEvent); ok {
			return c.SendTempMessage(msg.GroupUin, msg.SenderUin, m)
		}
	case MessageEventType:
		// 断言为统一消息事件
		if msg, ok := event.(MessageEvent); ok {
			// 通过tag来判断消息类型
			if Contains(msg.EventTags, "private_message") {
				return c.SendPrivateMessage(msg.SenderUin, m)
			} else if Contains(msg.EventTags, "group_message") {
				return c.SendGroupMessage(msg.GroupUin, m)
			} else if Contains(msg.EventTags, "temp_message") {
				return c.SendTempMessage(msg.GroupUin, msg.SenderUin, m)
			}
		}
	default:
		Error("发送消息时传入了不支持的消息事件")
	}
	return false, 0
}

func (c *CryoClient) Reply(event CryoMessageEvent, args ...interface{}) (ok bool, messageId uint32) {
	// 处理消息内容
	m := BuildMessage().Reply(event).Add(*ProcessMessageContent(args...))
	// 根据传入的事件来发送消息
	switch event.Type() {
	case PrivateMessageEventType:
		// 断言为私聊消息事件
		if msg, ok := event.(PrivateMessageEvent); ok {
			return c.SendPrivateMessage(msg.SenderUin, m.Reply(msg))
		}
	case GroupMessageEventType:
		// 断言为群消息事件
		if msg, ok := event.(GroupMessageEvent); ok {
			return c.SendGroupMessage(msg.GroupUin, m.Reply(msg))
		}
	case TempMessageEventType:
		// 断言为临时消息事件
		if msg, ok := event.(TempMessageEvent); ok {
			return c.SendTempMessage(msg.GroupUin, msg.SenderUin, m.Reply(msg))
		}
	case MessageEventType:
		// 断言为统一消息事件
		if msg, ok := event.(MessageEvent); ok {
			// 通过tag来判断消息类型
			if Contains(msg.EventTags, "private_message") {
				return c.SendPrivateMessage(msg.SenderUin, m.Reply(msg))
			} else if Contains(msg.EventTags, "group_message") {
				return c.SendGroupMessage(msg.GroupUin, m.Reply(msg))
			} else if Contains(msg.EventTags, "temp_message") {
				return c.SendTempMessage(msg.GroupUin, msg.SenderUin, m.Reply(msg))
			}
		}
	default:
		Error("发送消息时传入了不支持的消息事件")
	}
	return false, 0
}
