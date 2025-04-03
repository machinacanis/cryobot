/*

event.go

cryobot的事件结构体封装，事件的基本结构来自于LagrangGo，稍微修改了一下，让它风格更统一

*/

package event

import (
	"github.com/LagrangeDev/LagrangeGo/message"
	"github.com/go-json-experiment/json"
)

type CryoEventType uint32

const (
	PrivateMessageEventType                 CryoEventType = iota // 私聊消息事件类型
	GroupMessageEventType                                        // 群消息事件类型
	TempMessageEventType                                         // 临时消息事件类型
	NewFriendRequestEventType                                    // 新好友请求事件类型
	NewFriendEventType                                           // 新好友事件类型
	FriendRecallEventType                                        // 好友撤回事件类型
	FriendRenameEventType                                        // 好友改名事件类型
	FriendPokeEventType                                          // 好友戳一戳事件类型
	GroupMemberPermissionUpdatedEventType                        // 群成员权限变更事件类型
	GroupNameUpdatedEventType                                    // 群名称变更事件类型
	GroupMuteEventType                                           // 群禁言事件类型
	GroupRecallEventType                                         // 群撤回事件类型
	GroupMemberJoinRequestEventType                              // 群成员入群请求事件类型
	GroupMemberIncreaseEventType                                 // 群成员增加事件类型
	GroupMemberDecreaseEventType                                 // 群成员减少事件类型
	GroupDigestEventType                                         // 群精华消息事件类型
	GroupReactionEventType                                       // 群消息表情事件类型
	GroupMemberSpecialTitleUpdatedEventType                      // 群成员特殊头衔变更事件类型
	GroupInviteEventType                                         // 加群邀请事件类型
	BotConnectedEventType                                        // 机器人连接事件类型
	BotDisconnectedEventType                                     // 机器人断开连接事件类型
	CustomEventType                                              // 自定义事件类型
)

type (
	CryoEvent interface {
		Type() CryoEventType
		ToJson() []byte
		ToJsonString() string
	}

	CryoMessageEvent interface {
		CryoEvent
		replyDetail() (uint32, uint32, uint32, []message.IMessageElement)
	}

	// CryoBaseEvent 是CryoBot的事件总线上的事件结构体
	CryoBaseEvent struct {
		EventType   uint32 // 事件类型，是一个枚举
		EventID     string // 事件ID
		BotId       string // 机器人ID
		BotNickname string // 机器人昵称
		BotUin      uint32 // 机器人Uin
		BotUid      string // 机器人Uid
		Platform    string // 机器人平台
		Summary     string // 事件摘要
		Time        uint32 // 事件发生的时间戳
	}

	// MessageEvent 是CryoBot的消息事件结构体
	MessageEvent struct {
		CryoBaseEvent
		MessageId      uint32 // 消息ID
		SenderUin      uint32 // 消息发送者的Uin
		SenderUid      string // 消息发送者的Uid
		SenderNickname string // 消息发送者的昵称
		SenderCardname string // 消息发送者的备注名
		IsSenderFriend bool   // 消息发送者是否是好友

		MessageElements CryoMessage // 消息元素
	}

	// PrivateMessageEvent 私聊消息事件
	PrivateMessageEvent struct {
		MessageEvent
		InternalId uint32 // 内部ID
		ClientSeq  uint32 // 客户端序列号
		TargetUin  uint32 // 目标Uin
	}
	// GroupMessageEvent 群消息事件
	GroupMessageEvent struct {
		MessageEvent
		InternalId uint32 // 内部ID
		GroupUin   uint32 // 群号
		GroupName  string // 群名称
	}
	// TempMessageEvent 临时消息事件
	TempMessageEvent struct {
		MessageEvent
		InternalId    uint32 // 内部ID
		FromGroupUin  uint32 // 群号
		FromGroupName string // 群名称
	}
	// NewFriendRequestEvent 新好友请求事件
	NewFriendRequestEvent struct {
		CryoBaseEvent
		Uin      uint32
		Uid      string
		Nickname string
		Message  string
		From     string
	}
	// NewFriendEvent 新好友事件
	NewFriendEvent struct {
		CryoBaseEvent
		Uin      uint32
		Uid      string
		Nickname string
		Message  string
	}
	// FriendRecallEvent 好友撤回事件
	FriendRecallEvent struct {
		CryoBaseEvent
		Uin     uint32
		Uid     string
		Seqence uint64
		Random  uint32
	}
	// FriendRenameEvent 好友改名事件
	FriendRenameEvent struct {
		CryoBaseEvent
		IsSelf   bool
		Uin      uint32
		Uid      string
		Nickname string
	}
	// FriendPokeEvent 好友戳一戳事件
	FriendPokeEvent struct {
		CryoBaseEvent
		SenderUin uint32
		TargetUin uint32
		Suffix    string
		Action    string
	}
	// GroupMemberPermissionUpdatedEvent 群成员权限变更事件
	GroupMemberPermissionUpdatedEvent struct {
		CryoBaseEvent
		GroupUin uint32
		Uin      uint32
		Uid      string
		IsAdmin  bool
	}
	// GroupNameUpdatedEvent 群名称变更事件
	GroupNameUpdatedEvent struct {
		CryoBaseEvent
		GroupUin uint32
		Uin      uint32
		Uid      string
		NewName  string
	}
	// GroupMuteEvent 群禁言事件
	GroupMuteEvent struct {
		CryoBaseEvent
		GroupUin    uint32
		OperatorUin uint32
		OperatorUid string
		TargetUin   uint32
		TargetUid   string
		Duration    uint32
		isMuteAll   bool
	}
	// GroupRecallEvent 群撤回事件
	GroupRecallEvent struct {
		CryoBaseEvent
		GroupUin    uint32
		OperatorUin uint32
		OperatorUid string
		SenderUin   uint32
		SenderUid   string
		Seqence     uint64
		Random      uint32
	}
	// GroupMemberJoinRequestEvent 群成员入群请求事件
	GroupMemberJoinRequestEvent struct {
		CryoBaseEvent
		GroupUin       uint32
		SenderUin      uint32
		SenderUid      string
		SenderNickname string
		InviterUin     uint32
		InviterUid     string
		Answer         string
		RequestSeqence uint64
	}
	// GroupMemberIncreaseEvent 群成员增加事件
	GroupMemberIncreaseEvent struct {
		CryoBaseEvent
		GroupUin   uint32
		Uin        uint32
		Uid        string
		InviterUin uint32
		InviterUid string
		IsSelf     bool
	}
	// GroupMemberDecreaseEvent 群成员减少事件
	GroupMemberDecreaseEvent struct {
		CryoBaseEvent
		GroupUin uint32
		Uin      uint32
		Uid      string
		IsSelf   bool
	}
	// GroupDigestEvent 群精华消息事件
	GroupDigestEvent struct {
		CryoBaseEvent
		GroupUin         uint32
		MessageId        string
		InternalId       uint32
		SenderUin        uint32
		SenderUid        string
		SenderNickname   string
		OperatorUin      uint32
		OperatorUid      string
		OperatorNickname string
		IsRemove         bool
	}
	// GroupReactionEvent 群消息表态事件
	GroupReactionEvent struct {
		CryoBaseEvent
		GroupUin  uint32
		Uin       uint32
		Uid       string
		TargetSeq uint32
		IsAdd     bool
		IsEmoji   bool
		Code      string
		Count     uint32
	}
	// GroupMemberSpecialTitleUpdated 群成员特殊头衔变更事件
	GroupMemberSpecialTitleUpdated struct {
		CryoBaseEvent
		GroupUin uint32
		Uin      uint32
		Uid      string
		NewTitle string
	}
	// GroupInviteEvent 加群邀请事件
	GroupInviteEvent struct {
		CryoBaseEvent
		GroupUin        uint32
		GroupName       string
		InviterUin      uint32
		InviterUid      string
		InviterNickname string
		RequestSeqence  uint64
	}
	// BotConnectedEvent 机器人连接事件
	BotConnectedEvent struct {
		CryoBaseEvent
		Version string
	}
	// BotDisconnectedEvent 机器人断开连接事件
	BotDisconnectedEvent struct {
		CryoBaseEvent
	}
	CustomEvent struct {
		CryoBaseEvent
		summury string      // 摘要
		payload interface{} // 负载
	}
)

func (e PrivateMessageEvent) Type() CryoEventType {
	return PrivateMessageEventType
}

func (e GroupMessageEvent) Type() CryoEventType {
	return GroupMessageEventType
}

func (e TempMessageEvent) Type() CryoEventType {
	return TempMessageEventType
}

func (e NewFriendRequestEvent) Type() CryoEventType {
	return NewFriendRequestEventType
}

func (e NewFriendEvent) Type() CryoEventType {
	return NewFriendEventType
}

func (e FriendRecallEvent) Type() CryoEventType {
	return FriendRecallEventType
}

func (e FriendRenameEvent) Type() CryoEventType {
	return FriendRenameEventType
}

func (e FriendPokeEvent) Type() CryoEventType {
	return FriendPokeEventType
}

func (e GroupMemberPermissionUpdatedEvent) Type() CryoEventType {
	return GroupMemberPermissionUpdatedEventType
}

func (e GroupNameUpdatedEvent) Type() CryoEventType {
	return GroupNameUpdatedEventType
}

func (e GroupMuteEvent) Type() CryoEventType {
	return GroupMuteEventType
}

func (e GroupRecallEvent) Type() CryoEventType {
	return GroupRecallEventType
}

func (e GroupMemberJoinRequestEvent) Type() CryoEventType {
	return GroupMemberJoinRequestEventType
}

func (e GroupMemberIncreaseEvent) Type() CryoEventType {
	return GroupMemberIncreaseEventType
}

func (e GroupMemberDecreaseEvent) Type() CryoEventType {
	return GroupMemberDecreaseEventType
}

func (e GroupDigestEvent) Type() CryoEventType {
	return GroupDigestEventType
}

func (e GroupReactionEvent) Type() CryoEventType {
	return GroupReactionEventType
}

func (e GroupMemberSpecialTitleUpdated) Type() CryoEventType {
	return GroupMemberSpecialTitleUpdatedEventType
}

func (e GroupInviteEvent) Type() CryoEventType {
	return GroupInviteEventType
}

func (e BotConnectedEvent) Type() CryoEventType {
	return BotConnectedEventType
}

func (e BotDisconnectedEvent) Type() CryoEventType {
	return BotDisconnectedEventType
}

func (e CustomEvent) Type() CryoEventType {
	return CustomEventType
}

func (e PrivateMessageEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupMessageEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e TempMessageEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e NewFriendRequestEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e NewFriendEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e FriendRecallEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e FriendRenameEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e FriendPokeEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupMemberPermissionUpdatedEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupNameUpdatedEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupMuteEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupRecallEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupMemberJoinRequestEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupMemberIncreaseEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupMemberDecreaseEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupDigestEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupReactionEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupMemberSpecialTitleUpdated) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e GroupInviteEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e BotConnectedEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e BotDisconnectedEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e CustomEvent) ToJson() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

func (e PrivateMessageEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupMessageEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e TempMessageEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e NewFriendRequestEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e NewFriendEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e FriendRecallEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e FriendRenameEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e FriendPokeEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupMemberPermissionUpdatedEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupNameUpdatedEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupMuteEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupRecallEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupMemberJoinRequestEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupMemberIncreaseEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupMemberDecreaseEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupDigestEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupReactionEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupMemberSpecialTitleUpdated) ToJsonString() string {
	return string(e.ToJson())
}

func (e GroupInviteEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e BotConnectedEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e BotDisconnectedEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e CustomEvent) ToJsonString() string {
	return string(e.ToJson())
}

func (e PrivateMessageEvent) replyDetail() (uint32, uint32, uint32, []message.IMessageElement) {
	return e.MessageId, e.SenderUin, e.Time, e.MessageElements.ToIMessageElements()
}

func (e GroupMessageEvent) replyDetail() (uint32, uint32, uint32, []message.IMessageElement) {
	return e.MessageId, e.SenderUin, e.Time, e.MessageElements.ToIMessageElements()
}
