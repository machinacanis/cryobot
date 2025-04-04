package cryobot

import (
	"github.com/go-json-experiment/json"
	"os"
)

type CryoClientInfo struct {
	Id        string `json:"id"`
	Platform  string `json:"Platform"`
	Version   string `json:"Version"`
	DeviceNum int    `json:"device_num"`
	Signature string `json:"signature"`
	Uin       int    `json:"uin"`
	Uid       string `json:"uid"`
}

func ReadClientInfos() ([]CryoClientInfo, error) {
	data, err := os.ReadFile("client_infos.json")
	if err != nil {
		return nil, err
	}

	var clientInfos []CryoClientInfo
	err = json.Unmarshal(data, &clientInfos)
	if err != nil {
		return nil, err
	}

	return clientInfos, nil
}

func WriteClientInfos(clientInfos []CryoClientInfo) error {
	data, err := json.Marshal(clientInfos)
	if err != nil {
		return err
	}

	err = os.WriteFile("client_infos.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func SaveClientInfo(clientInfo CryoClientInfo) error {
	// 首先尝试读取现有的客户端信息
	clientInfos, err := ReadClientInfos()
	if err != nil { // 读取失败，可能是文件不存在
		if os.IsNotExist(err) {
			// 文件不存在，创建一个新的切片
			clientInfos = []CryoClientInfo{}
		} else {
			return err // 其他错误
		}
	}
	// 检测是否已经存在botid相同的客户端信息
	updateFlag := false
	for _, info := range clientInfos {
		if info.Id == clientInfo.Id {
			// 如果存在，则更新该信息
			updateFlag = true
			break
		}
	}
	if updateFlag {
		// 如果存在，则更新该信息
		for i, info := range clientInfos {
			if info.Id == clientInfo.Id {
				clientInfos[i] = clientInfo
				break
			}
		}
	} else {
		// 如果不存在，则添加新的信息
		clientInfos = append(clientInfos, clientInfo)
	}
	err = WriteClientInfos(clientInfos)
	return err
}

func RemoveClientInfo(botId string) error {
	// 读取现有的客户端信息
	clientInfos, err := ReadClientInfos()
	if err != nil {
		return err // 读取失败，可能是文件不存在
	}

	// 创建一个新的切片来存储更新后的客户端信息
	var updatedClientInfos []CryoClientInfo

	// 遍历现有的客户端信息，排除要删除的项
	for _, info := range clientInfos {
		if info.Id != botId {
			updatedClientInfos = append(updatedClientInfos, info)
		}
	}

	// 将更新后的客户端信息写回文件
	err = WriteClientInfos(updatedClientInfos)
	return err
}
