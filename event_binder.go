package cryobot

import (
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
	uuid "github.com/satori/go.uuid"
	"time"
)

// EventBind 绑定LagrangeGo的事件到cryobot的事件总线
func EventBind(cc *CryoClient) {

	Infof("%s[Cryo] 正在将 %d 的消息事件绑定到事件总线", lavender, cc.Client.Uin)
	// 断开连接
	cc.Client.DisconnectedEvent.Subscribe(func(client *client.QQClient, event *client.DisconnectedEvent) {
		SendBotConnectedEvent(cc)
	})

	// 私聊消息
	cc.Client.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {
		PublishAsync(PrivateMessageEvent{
			MessageEvent: MessageEvent{
				BaseEvent: BaseEvent{
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
				GroupUin:        event.Sender.Uin,
				GroupName:       event.Sender.Nickname,
			},
			InternalId: event.InternalID,
			ClientSeq:  event.ClientSeq,
			TargetUin:  event.Target,
		})
	})

	// 群聊消息
	cc.Client.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		PublishAsync(GroupMessageEvent{
			MessageEvent: MessageEvent{
				BaseEvent: BaseEvent{
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
				GroupUin:        event.GroupUin,
				GroupName:       event.GroupName,
			},
			InternalId: event.InternalID,
		})
	})

	cc.Client.TempMessageEvent.Subscribe(func(client *client.QQClient, event *message.TempMessage) {
		PublishAsync(TempMessageEvent{
			MessageEvent: MessageEvent{
				BaseEvent: BaseEvent{
					EventType:   uint32(TempMessageEventType),
					EventId:     uuid.NewV4().String(),
					EventTags:   []string{"temp_message", "message"},
					BotId:       cc.Id,
					BotNickname: cc.Nickname,
					BotUin:      uint32(cc.Uin),
					BotUid:      cc.Uid,
					Platform:    cc.Platform,
					Summary:     "TempMessageEvent",
					Time:        uint32(time.Now().Unix()),
				},
				MessageId:       event.ID,
				SenderUin:       event.Sender.Uin,
				SenderUid:       event.Sender.UID,
				SenderNickname:  event.Sender.Nickname,
				SenderCardname:  event.Sender.CardName,
				IsSenderFriend:  event.Sender.IsFriend,
				MessageElements: *FromLagrangeMessage(event.Elements),
				GroupUin:        event.GroupUin,
				GroupName:       event.GroupName,
			},
		})
	})

	Infof("%s[Cryo] %d 的消息事件绑定完成", lavender, cc.Client.Uin)
}
