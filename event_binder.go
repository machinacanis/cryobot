package cryobot

import (
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
	uuid "github.com/satori/go.uuid"
)

// EventBind 绑定LagrangeGo的事件到cryobot的事件总线
func EventBind(cc *CryoClient) {

	Infof("%s[Cryo] 正在将 %d 的消息事件绑定到事件总线", lavender, cc.LagrangeClient.Uin)
	// 断开连接
	cc.LagrangeClient.DisconnectedEvent.Subscribe(func(client *client.QQClient, event *client.DisconnectedEvent) {
		SendBotConnectedEvent(cc)
	})

	// 私聊消息
	cc.LagrangeClient.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {
		PublishAsync(PrivateMessageEvent{
			MessageEvent: MessageEvent{
				CryoBaseEvent: CryoBaseEvent{
					EventType:   uint32(PrivateMessageEventType),
					EventId:     uuid.NewV4().String(),
					EventTags:   []string{"private_message", "message"},
					BotId:       cc.Id,
					BotNickname: cc.Nickname,
					BotUin:      uint32(cc.Uin),
					BotUid:      cc.Uid,
					Platform:    cc.Platform,
					Summary:     "PrivateMessageEvent",
					Time:        event.Time,
				},
				MessageId:       event.ID,
				SenderUin:       event.Sender.Uin,
				SenderUid:       event.Sender.UID,
				SenderNickname:  event.Sender.Nickname,
				SenderCardname:  event.Sender.CardName,
				IsSenderFriend:  event.Sender.IsFriend,
				MessageElements: *FromLagrangeMessage(event.Elements),
			},
			InternalId: event.InternalID,
			ClientSeq:  event.ClientSeq,
			TargetUin:  event.Target,
		})
	})

	// 群聊消息
	cc.LagrangeClient.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		PublishAsync(GroupMessageEvent{
			MessageEvent: MessageEvent{
				CryoBaseEvent: CryoBaseEvent{
					EventType:   uint32(GroupMessageEventType),
					EventId:     uuid.NewV4().String(),
					EventTags:   []string{"group_message", "message"},
					BotId:       cc.Id,
					BotNickname: cc.Nickname,
					BotUin:      uint32(cc.Uin),
					BotUid:      cc.Uid,
					Platform:    cc.Platform,
					Summary:     "GroupMessageEvent",
					Time:        event.Time,
				},
				MessageId:       event.ID,
				SenderUin:       event.Sender.Uin,
				SenderUid:       event.Sender.UID,
				SenderNickname:  event.Sender.Nickname,
				SenderCardname:  event.Sender.CardName,
				IsSenderFriend:  event.Sender.IsFriend,
				MessageElements: *FromLagrangeMessage(event.Elements),
			},
			InternalId: event.InternalID,
			GroupUin:   event.GroupUin,
			GroupName:  event.GroupName,
		})
	})

	Infof("%s[Cryo] %d 的消息事件绑定完成", lavender, cc.LagrangeClient.Uin)
}
