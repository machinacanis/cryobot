package cryobot

import (
	"github.com/go-json-experiment/json"
	"github.com/sirupsen/logrus"
	"os"
)

var conf Config
var DefaultSignServer = "https://sign.lagrangecore.org/api/sign/30366"

type Config struct {
	LogLevel                     logrus.Level
	LogFormat                    *logrus.Formatter
	SignServers                  []string `json:"sign_servers,omitempty,omitzero"`                    // 签名服务器列表
	EnableClientAutoSave         bool     `json:"enable_client_save,omitempty,omitzero"`              // 是否启用客户端信息自动保存
	EnablePrintLogo              bool     `json:"enable_print_logo,omitempty,omitzero"`               // 是否启用logo打印
	EnableConnectPrintMiddleware bool     `json:"enable_connect_print_middleware,omitempty,omitzero"` // 是否启用内置的Bot连接打印中间件
	EnableMessagePrintMiddleware bool     `json:"enable_message_print_middleware,omitempty,omitzero"` // 是否启用内置的消息打印中间件
	EnableEventDebugMiddleware   bool     `json:"enable_event_debug_middleware,omitempty,omitzero"`   // 是否启用内置的事件调试中间件
}

func ReadCryoConfig() (Config, error) {
	data, err := os.ReadFile("cryobot_config.json")
	c := Config{}
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
func WriteCryoConfig(config Config) error {
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
