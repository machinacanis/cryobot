/*

lagrange_client.go

*/

package client

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/auth"
	"github.com/machinacanis/cryobot/config"
	"github.com/machinacanis/cryobot/log"
	"github.com/machinacanis/cryobot/utils"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var connectedClients []*LagrangeClient // 连接的客户端列表

// LagrangeClient LagrangeGo库的客户端封装，封装后可以实现更方便的多bot管理
type LagrangeClient struct {
	BotId     string
	Client    *client.QQClient
	Platform  string
	Version   string
	DeviceUin int
	Uin       int
	Uid       string
}

// New 初始化一个新的Lagrange客户端
func (lc *LagrangeClient) New() {
	lc.BotId = utils.NewUUID() // 给Bot客户端分配一个唯一的UUID

	// 默认平台和版本
	if lc.Platform == "" {
		lc.Platform = "linux"
	}
	if lc.Version == "" {
		lc.Version = "3.2.15-30366"
	}

	appInfo := auth.AppList[lc.Platform][lc.Version]
	qqClientInstance := client.NewClient(0, "")
	qqClientInstance.SetLogger(pLogger) // 替换日志记录器，详见client/protocol_logger.go以及log/logger.go
	qqClientInstance.UseVersion(appInfo)
	qqClientInstance.AddSignServer(config.Conf.SignServers...)
	// 随机一个设备号
	lc.DeviceUin = utils.RandomDeviceNumber()
	qqClientInstance.UseDevice(auth.NewDeviceInfo(lc.DeviceUin))
	/*
		新建客户端时不需要读取签名文件

		data, err := os.ReadFile(fmt.Sprintf("%s.bin", lc.BotId))
		if err != nil {
			log.Warn("未找到签名文件：", err)
		} else {
			sig, err := auth.UnmarshalSigInfo(data, true)
			if err != nil {
				log.Warn("读取签名文件失败：", err)
			} else {
				qqClientInstance.UseSig(sig)
			}
		}
	*/
	lc.Client = qqClientInstance // 将LagrangeGo的客户端实例赋值给LagrangeClient
}

// GetSignature 获取当前客户端的签名信息
func (lc *LagrangeClient) GetSignature() string {
	data, err := lc.Client.Sig().Marshal()
	if err != nil {
		log.Error("序列化签名时出现错误：", err)
		return ""
	}
	// 将二进制的签名直接编码到字符串
	sig := base64.StdEncoding.EncodeToString(data)
	return sig
}

// UseSignature 使用指定的签名信息
func (lc *LagrangeClient) UseSignature(sig string) {
	// 将字符串解码为二进制
	data, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		log.Error("解码签名时出现错误：", err)
		return
	}
	// 反序列化签名
	sigInfo, err := auth.UnmarshalSigInfo(data, true)
	if err != nil {
		log.Error("反序列化签名时出现错误：", err)
		return
	}
	lc.Client.UseSig(sigInfo)
}

// Open 尝试重建已经存在的Lagrange客户端
func (lc *LagrangeClient) Open(botId string) error {
	// 尝试在客户端信息文件中查找有相同id的客户端信息
	var sig string
	clientInfos, err := ReadClientInfos()
	if err != nil {
		log.Error("读取客户端信息时出现错误：", err)
		return err
	}
	foundFlag := false
	for _, info := range clientInfos {
		if info.BotId == botId {
			lc.BotId = info.BotId
			lc.Platform = info.Platform
			lc.Version = info.Version
			lc.DeviceUin = info.DeviceUin
			sig = info.Signature
			lc.Uin = info.Uin
			lc.Uid = info.Uid
			foundFlag = true
			break
		}
	}
	if lc.BotId == "" || !foundFlag { // 如果没有找到相应的客户端信息，则返回异常
		log.Error("未找到匹配的客户端信息")
		return errors.New("未找到匹配的客户端信息")
	}
	appInfo := auth.AppList[lc.Platform][lc.Version]
	qqClientInstance := client.NewClient(0, "")
	qqClientInstance.SetLogger(pLogger) // 替换日志记录器，详见client/protocol_logger.go以及log/logger.go
	qqClientInstance.UseVersion(appInfo)
	qqClientInstance.AddSignServer(config.Conf.SignServers...)
	qqClientInstance.UseDevice(auth.NewDeviceInfo(lc.DeviceUin))
	lc.Client = qqClientInstance
	lc.UseSignature(sig) // 使用签名信息
	return nil
}

// Save 将当前客户端的信息保存到文件中
func (lc *LagrangeClient) Save() error {
	clientInfo := LagrangeClientInfo{
		BotId:     lc.BotId,
		Signature: lc.GetSignature(),
		Platform:  lc.Platform,
		Version:   lc.Version,
		DeviceUin: lc.DeviceUin,
		Uin:       lc.Uin,
		Uid:       lc.Uid,
	}
	err := SaveClientInfo(clientInfo)
	return err
}

// Login 执行登录操作
func (lc *LagrangeClient) Login() error {
	// 自动登录逻辑
	sig := lc.Client.Sig()
	if sig != nil {
		err := lc.Client.FastLogin()
		if err != nil {
			log.Warn("快速登录失败，切换到二维码登录")
		} else {
			// 登录成功，获取签名里的信息
			lc.Uin = int(lc.Client.Sig().Uin)
			lc.Uid = lc.Client.Sig().UID
			log.Infof("%d 已成功连接", lc.Uin)
			lc.Client.DisconnectedEvent.Subscribe(func(client *client.QQClient, event *client.DisconnectedEvent) { // 订阅断开连接事件
				log.Infof("%d 连接已断开：%v", lc.Uin, event.Message)
			})
			err = lc.Save()
			if err != nil {
				return err
			} // 保存登录信息
			return nil
		}
	}
	// 如果没有签名或者快速登录失败，则使用二维码登录
	//获取二维码
	png, _, err := lc.Client.FetchQRCodeDefault()
	if err != nil {
		log.Error("登录时出现错误：", err)
		return err
	}
	//保存本地二维码
	qrcodePath := fmt.Sprintf("QRCode_%s.png", lc.BotId)
	err = os.WriteFile(qrcodePath, png, 0644)
	if err != nil {
		log.Error("写入二维码图片时出现错误：", err)
		return err
	}
	log.Infof("登录二维码已保存到 %s", qrcodePath)
	//轮询登录状态
	for {
		retCode, err := lc.Client.GetQRCodeResult()
		if err != nil {
			logrus.Errorln(err)
			return err
		}
		// 等待扫码
		if retCode.Waitable() {
			time.Sleep(3 * time.Second)
			continue
		}
		if !retCode.Success() {
			return errors.New(retCode.Name())
		}
		break
	}
	_, err = lc.Client.QRCodeLogin()
	if err != nil {
		log.Error("login err:", err)
		return err
	}
	// 登录成功，获取签名里的信息
	lc.Uin = int(lc.Client.Sig().Uin)
	lc.Uid = lc.Client.Sig().UID
	log.Infof("%d 已成功连接", lc.Uin)
	lc.Client.DisconnectedEvent.Subscribe(func(client *client.QQClient, event *client.DisconnectedEvent) { // 订阅断开连接事件
		log.Infof("%d 连接已断开：%v", lc.Uin, event.Message)
	})
	err = lc.Save()
	if err != nil {
		return err
	} // 保存登录信息
	return nil
}

// Logout 执行登出操作
func (lc *LagrangeClient) Logout() error {
	lc.Client.Disconnect()
	log.Infof("%d 已成功登出", lc.Uin)
	return nil
}

// Connect 尝试查询并连接到指定的bot客户端
func Connect(botId string) error {
	lc := &LagrangeClient{}
	err := lc.Open(botId)
	if err != nil {
		return err
	}
	log.Infof("正在连接Bot：%s (%d)", lc.BotId, lc.Uin)
	err = lc.Login()
	if err != nil {
		return err
	}
	connectedClients = append(connectedClients, lc)
	return nil
}

// ConnectAll 尝试连接所有已保存的bot客户端
func ConnectAll() {
	clientInfos, err := ReadClientInfos()
	if err != nil {
		log.Error("读取Bot信息时出现错误：", err)
		return
	}
	if len(clientInfos) == 0 {
		log.Info("没有找到Bot信息")
		return
	}
	for _, info := range clientInfos {
		err := Connect(info.BotId)
		if err != nil {
			log.Error("连接Bot时出现错误：", err)
			log.Error("已自动清除失效的客户端信息，请重新登录")
		}
	}
}

func ConnectNew() {
	lc := &LagrangeClient{}
	lc.New()
	err := lc.Login()
	if err != nil {
		log.Error("登录时出现错误：", err)
		return
	}
	connectedClients = append(connectedClients, lc)
}
