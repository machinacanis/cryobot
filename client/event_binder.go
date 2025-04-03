/*

event_binder.go

*/

package client

import (
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
	cryoevent "github.com/machinacanis/cryobot/event"
	"github.com/machinacanis/cryobot/log"
	uuid "github.com/satori/go.uuid"
)

// LagrangeEventBind 绑定LagrangeGo的事件到cryobot的事件总线
func LagrangeEventBind(lc *LagrangeClient) {

	log.Debugf("正在将 %d 的消息事件绑定到事件总线", lc.Client.Uin)
	// 断开连接
	lc.Client.DisconnectedEvent.Subscribe(func(client *client.QQClient, event *client.DisconnectedEvent) {
		SendBotConnectedEvent(lc)
	})

	// 私聊消息
	lc.Client.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {
		AsyncPublish(cryoevent.PrivateMessageEvent{
			MessageEvent: cryoevent.MessageEvent{
				CryoBaseEvent: cryoevent.CryoBaseEvent{
					EventType:   uint32(cryoevent.PrivateMessageEventType),
					EventID:     uuid.NewV4().String(),
					BotId:       lc.BotId,
					BotNickname: lc.Nickname,
					BotUin:      uint32(lc.Uin),
					BotUid:      lc.Uid,
					Platform:    lc.Platform,
					Summary:     "PrivateMessageEvent",
					Time:        event.Time,
				},
				MessageId:       event.ID,
				SenderUin:       event.Sender.Uin,
				SenderUid:       event.Sender.UID,
				SenderNickname:  event.Sender.Nickname,
				SenderCardname:  event.Sender.CardName,
				IsSenderFriend:  event.Sender.IsFriend,
				MessageElements: *cryoevent.FromLagrangeMessage(event.Elements),
			},
			InternalId: event.InternalID,
			ClientSeq:  event.ClientSeq,
			TargetUin:  event.Target,
		})
	})

	// 群聊消息
	lc.Client.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		AsyncPublish(cryoevent.GroupMessageEvent{
			MessageEvent: cryoevent.MessageEvent{
				CryoBaseEvent: cryoevent.CryoBaseEvent{
					EventType:   uint32(cryoevent.GroupMessageEventType),
					EventID:     uuid.NewV4().String(),
					BotId:       lc.BotId,
					BotNickname: lc.Nickname,
					BotUin:      uint32(lc.Uin),
					BotUid:      lc.Uid,
					Platform:    lc.Platform,
					Summary:     "GroupMessageEvent",
					Time:        event.Time,
				},
				MessageId:       event.ID,
				SenderUin:       event.Sender.Uin,
				SenderUid:       event.Sender.UID,
				SenderNickname:  event.Sender.Nickname,
				SenderCardname:  event.Sender.CardName,
				IsSenderFriend:  event.Sender.IsFriend,
				MessageElements: *cryoevent.FromLagrangeMessage(event.Elements),
			},
			InternalId: event.InternalID,
			GroupUin:   event.GroupUin,
			GroupName:  event.GroupName,
		})
	})

	log.Debugf("%d 的消息事件绑定完成", lc.Client.Uin)
}
