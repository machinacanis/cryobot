package cryobot

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

func SendBotConnectedEvent(cc *CryoClient) {
	PublishAsync(BotConnectedEvent{
		CryoBaseEvent: CryoBaseEvent{
			EventType:   uint32(BotConnectedEventType),
			EventId:     uuid.NewV4().String(),
			EventTags:   []string{"system", "bot"},
			BotId:       cc.Id,
			BotNickname: cc.Nickname,
			BotUin:      uint32(cc.Uin),
			BotUid:      cc.Uid,
			Platform:    cc.Platform,
			Summary:     "BotConnectedEvent",
			Time:        uint32(time.Now().Unix()),
		},
		Version: cc.Version,
	})
}

func SendBotDisconnectedEvent(cc *CryoClient) {
	PublishAsync(BotDisconnectedEvent{
		CryoBaseEvent: CryoBaseEvent{
			EventType:   uint32(BotDisconnectedEventType),
			EventId:     uuid.NewV4().String(),
			EventTags:   []string{"system", "bot"},
			BotId:       cc.Id,
			BotNickname: cc.Nickname,
			BotUin:      uint32(cc.Uin),
			BotUid:      cc.Uid,
			Platform:    cc.Platform,
			Summary:     "BotDisconnectedEvent",
			Time:        uint32(time.Now().Unix()),
		},
	})
}
