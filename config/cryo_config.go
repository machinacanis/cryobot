package config

import (
	"github.com/go-json-experiment/json"
	"github.com/sirupsen/logrus"
	"os"
)

var Conf CryoConfig
var DefaultSignServer = "https://sign.lagrangecore.org/api/sign/30366"

type CryoConfig struct {
	LogLevel                     logrus.Level
	SignServers                  []string `json:"sign_servers,omitempty,omitzero"`                    // 签名服务器列表
	EnableConnectPrintMiddleware bool     `json:"enable_connect_print_middleware,omitempty,omitzero"` // 是否启用内置的Bot连接打印中间件
	EnableMessagePrintMiddleware bool     `json:"enable_message_print_middleware,omitempty,omitzero"` // 是否启用内置的消息打印中间件
	EnableEventDebugMiddleware   bool     `json:"enable_event_debug_middleware,omitempty,omitzero"`   // 是否启用内置的事件调试中间件
}

func ReadCryoConfig() (CryoConfig, error) {
	data, err := os.ReadFile("cryobot_config.json")
	c := CryoConfig{}
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	return c, nil
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
