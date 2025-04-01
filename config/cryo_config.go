package config

import (
	"github.com/go-json-experiment/json"
	"github.com/machinacanis/cryobot/log"
	"os"
)

var Conf CryoConfig
var defaultSignServer = "https://sign.lagrangecore.org/api/sign/30366"

func init() {
	err := ReadCryoConfig()
	if err != nil {
		// 如果读取配置文件失败，使用默认配置
		Conf = CryoConfig{
			SignServers:   []string{defaultSignServer},
			EnableBackend: false,
		}
		err = WriteCryoConfig(Conf)
		if err != nil {
			log.Error("无法创建配置文件：", err)
		} else {
			log.Info("配置文件不存在，已创建默认配置文件")
		}
	} else {
		if len(Conf.SignServers) == 0 {
			Conf.SignServers = []string{defaultSignServer}
			err = WriteCryoConfig(Conf)
			if err != nil {
				log.Error("无法更新配置文件：", err)
			} else {
				log.Info("配置文件已更新，使用默认签名服务器")
			}
		}
	}
}

type CryoConfig struct {
	SignServers                []string `json:"sign_servers"`                 // 签名服务器列表
	EnableBackend              bool     `json:"enable_backend"`               // 是否启用后端
	EnableMessageDeduplication bool     `json:"enable_message_deduplication"` // 是否启用消息去重
	EnableMessagePrint         bool     `json:"enable_message_print"`         // 是否启用消息打印
}

func ReadCryoConfig() error {
	data, err := os.ReadFile("cryobot_config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &Conf)
	if err != nil {
		return err
	}

	return nil
}
func WriteCryoConfig(config CryoConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile("cryobot_config.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
