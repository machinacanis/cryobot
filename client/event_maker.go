package client

import (
	"github.com/machinacanis/cryobot/event"
	uuid "github.com/satori/go.uuid"
	"time"
)

func SendBotConnectedEvent(lc *LagrangeClient) {
	AsyncPublish(event.BotConnectedEvent{
		CryoBaseEvent: event.CryoBaseEvent{
			EventType:   uint32(event.BotConnectedEventType),
			EventID:     uuid.NewV4().String(),
			BotId:       lc.BotId,
			BotNickname: lc.Nickname,
			BotUin:      uint32(lc.Uin),
			BotUid:      lc.Uid,
			Platform:    lc.Platform,
			Summary:     "BotConnectedEvent",
			Time:        uint32(time.Now().Unix()),
		},
		Version: lc.Version,
	})
}

func SendBotDisconnectedEvent(lc *LagrangeClient) {
	AsyncPublish(event.BotDisconnectedEvent{
		CryoBaseEvent: event.CryoBaseEvent{
			EventType:   uint32(event.BotDisconnectedEventType),
			EventID:     uuid.NewV4().String(),
			BotId:       lc.BotId,
			BotNickname: lc.Nickname,
			BotUin:      uint32(lc.Uin),
			BotUid:      lc.Uid,
			Platform:    lc.Platform,
			Summary:     "BotDisconnectedEvent",
			Time:        uint32(time.Now().Unix()),
		},
	})
}
