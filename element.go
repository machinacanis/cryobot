package cryobot

import (
	"fmt"
	lagrangeMessage "github.com/LagrangeDev/LagrangeGo/message"
	"io"
)

// 定义一系列LagrangeGo的类型别名
// 方便调用以及实现一些额外的方法

// ElementType 消息元素类型
type ElementType int

const (
	Text       ElementType = iota // 文本
	Image                         // 图片
	Face                          // 表情
	At                            // 艾特
	Reply                         // 回复
	Service                       // 服务
	Forward                       // 转发
	File                          // 文件
	Voice                         // 语音
	Video                         // 视频
	LightApp                      // 轻应用
	RedBag                        // 红包（无实际意义）
	MarketFace                    // 魔法表情
)

type Element interface {
	Type() ElementType
}

type TextElement struct {
	lagrangeMessage.TextElement
}
type AtElement struct {
	lagrangeMessage.AtElement
}
type FaceElement struct {
	lagrangeMessage.FaceElement
}
type ReplyElement struct {
	lagrangeMessage.ReplyElement
}
type VoiceElement struct {
	lagrangeMessage.VoiceElement
}
type ImageElement struct {
	lagrangeMessage.ImageElement
}
type FileElement struct {
	lagrangeMessage.FileElement
}
type ShortVideoElement struct {
	lagrangeMessage.ShortVideoElement
}
type LightAppElement struct {
	lagrangeMessage.LightAppElement
}
type XMLElement struct {
	lagrangeMessage.XMLElement
}
type ForwardMessageElement struct {
	lagrangeMessage.ForwardMessage
}
type MarketFaceElement struct {
	lagrangeMessage.MarketFaceElement
}

func (e *TextElement) Type() ElementType {
	return Text
}
func (e *AtElement) Type() ElementType {
	return At
}
func (e *FaceElement) Type() ElementType {
	return Face
}
func (e *ReplyElement) Type() ElementType {
	return Reply
}
func (e *VoiceElement) Type() ElementType {
	return Voice
}
func (e *ImageElement) Type() ElementType {
	return Image
}
func (e *FileElement) Type() ElementType {
	return File
}
func (e *ShortVideoElement) Type() ElementType {
	return Video
}
func (e *LightAppElement) Type() ElementType {
	return LightApp
}
func (e *XMLElement) Type() ElementType {
	return Service
}
func (e *ForwardMessageElement) Type() ElementType {
	return Forward
}
func (e *MarketFaceElement) Type() ElementType {
	return MarketFace
}

type CryoMessage struct {
	// 消息元素列表
	Elements []Element
}

// BuildMessage 构建一个新的CryoMessage实例
func BuildMessage(elements ...Element) *CryoMessage {
	// 构建一个新的CryoMessage实例
	return &CryoMessage{
		Elements: elements,
	}
}

// FromLagrangeMessage 将LagrangeGo的消息元素列表转换为CryoMessage实例
func FromLagrangeMessage(elements []lagrangeMessage.IMessageElement) *CryoMessage {
	// 将LagrangeGo的消息元素列表转换为CryoMessage实例
	var cryoElements []Element
	for _, element := range elements {
		switch e := element.(type) {
		case *lagrangeMessage.TextElement:
			cryoElements = append(cryoElements, &TextElement{*e})
		case *lagrangeMessage.AtElement:
			cryoElements = append(cryoElements, &AtElement{*e})
		case *lagrangeMessage.FaceElement:
			cryoElements = append(cryoElements, &FaceElement{*e})
		case *lagrangeMessage.ReplyElement:
			cryoElements = append(cryoElements, &ReplyElement{*e})
		case *lagrangeMessage.VoiceElement:
			cryoElements = append(cryoElements, &VoiceElement{*e})
		case *lagrangeMessage.ImageElement:
			cryoElements = append(cryoElements, &ImageElement{*e})
		case *lagrangeMessage.FileElement:
			cryoElements = append(cryoElements, &FileElement{*e})
		case *lagrangeMessage.ShortVideoElement:
			cryoElements = append(cryoElements, &ShortVideoElement{*e})
		case *lagrangeMessage.LightAppElement:
			cryoElements = append(cryoElements, &LightAppElement{*e})
		case *lagrangeMessage.XMLElement:
			cryoElements = append(cryoElements, &XMLElement{*e})
		case *lagrangeMessage.ForwardMessage:
			cryoElements = append(cryoElements, &ForwardMessageElement{*e})
		case *lagrangeMessage.MarketFaceElement:
			cryoElements = append(cryoElements, &MarketFaceElement{*e})
		default:
			continue
		}
	}
	return BuildMessage(cryoElements...)
}

func LagrangeMessageToString(elements []lagrangeMessage.IMessageElement) string {
	// 将LagrangeGo的消息元素列表转换为字符串
	var result string
	for _, element := range elements {
		switch e := element.(type) {
		case *lagrangeMessage.TextElement:
			result += e.Content
		case *lagrangeMessage.AtElement:
			result += fmt.Sprintf("@%d ", e.TargetUin)
		case *lagrangeMessage.FaceElement:
			result += fmt.Sprintf("[Face%d] ", e.FaceID)
		case *lagrangeMessage.ReplyElement:
			result += fmt.Sprintf("回复@(%d): %s", e.SenderUin, LagrangeMessageToString(e.Elements))
		case *lagrangeMessage.VoiceElement:
			result += "[语音]"
		case *lagrangeMessage.ImageElement:
			result += "[图像]"
		case *lagrangeMessage.FileElement:
			result += "[文件]"
		case *lagrangeMessage.ShortVideoElement:
			result += "[视频]"
		case *lagrangeMessage.LightAppElement:
			result += "[轻应用]"
		case *lagrangeMessage.XMLElement:
			result += "[服务]"
		case *lagrangeMessage.ForwardMessage:
			result += "[转发消息]"
		case *lagrangeMessage.MarketFaceElement:
			result += "[魔法表情]"
		default:
			continue
		}
	}
	return result
}

// Add 将一个CryoMessage添加到当前消息中
func (m *CryoMessage) Add(msg CryoMessage) *CryoMessage {
	// 将CryoMessage添加到当前消息中
	m.Elements = append(m.Elements, msg.Elements...)
	return m // 返回当前消息对象，以便链式调用
}

func (m *CryoMessage) ToIMessageElements() []lagrangeMessage.IMessageElement {
	m.Check() // 检查消息元素的合法性

	// 将CryoMessage转换为LagrangeGo的消息元素列表
	var elements []lagrangeMessage.IMessageElement
	for _, element := range m.Elements {
		switch e := element.(type) {
		case *TextElement:
			elements = append(elements, &e.TextElement)
		case *AtElement:
			elements = append(elements, &e.AtElement)
		case *FaceElement:
			elements = append(elements, &e.FaceElement)
		case *ReplyElement:
			elements = append(elements, &e.ReplyElement)
		case *VoiceElement:
			elements = append(elements, &e.VoiceElement)
		case *ImageElement:
			elements = append(elements, &e.ImageElement)
		case *FileElement:
			elements = append(elements, &e.FileElement)
		case *ShortVideoElement:
			elements = append(elements, &e.ShortVideoElement)
		case *LightAppElement:
			elements = append(elements, &e.LightAppElement)
		case *XMLElement:
			elements = append(elements, &e.XMLElement)
		case *ForwardMessageElement:
			elements = append(elements, &e.ForwardMessage)
		case *MarketFaceElement:
			elements = append(elements, &e.MarketFaceElement)
		default:
			continue
		}
	}
	return elements
}

func (m *CryoMessage) ToString() string {
	// 将CryoMessage转换为字符串
	var result string
	for _, element := range m.Elements {
		switch e := element.(type) {
		case *TextElement:
			result += e.TextElement.Content
		case *AtElement:
			result += fmt.Sprintf("@%d ", e.AtElement.TargetUin)
		case *FaceElement:
			result += fmt.Sprintf("[Face%d] ", e.FaceElement.FaceID)
		case *ReplyElement:
			result += fmt.Sprintf("回复@%d: %s", e.ReplyElement.SenderUin, LagrangeMessageToString(e.Elements))
		case *VoiceElement:
			result += "[语音]"
		case *ImageElement:
			result += "[图像]"
		case *FileElement:
			result += "[文件]"
		case *ShortVideoElement:
			result += "[视频]"
		case *LightAppElement:
			result += "[轻应用]"
		case *XMLElement:
			result += "[服务]"
		case *ForwardMessageElement:
			result += "[转发消息]"
		case *MarketFaceElement:
			result += "[魔法表情]"
		default:
			continue
		}
	}
	return result
}

func (m *CryoMessage) Check() {
	// Reply元素只能有一个，如果有多个，则只保留第一个
	replyCount := 0
	for i, element := range m.Elements {
		if _, ok := element.(*ReplyElement); ok {
			replyCount++
			if replyCount > 1 {
				m.Elements = append(m.Elements[:i], m.Elements[i+1:]...)
				i--
			}
		}
	}
	// Image元素最多20个，如果超过20个，则只保留前20个
	imageCount := 0
	for i, element := range m.Elements {
		if _, ok := element.(*ImageElement); ok {
			imageCount++
			if imageCount > 20 {
				m.Elements = append(m.Elements[:i], m.Elements[i+1:]...)
				i--
			}
		}
	}
}

func (m *CryoMessage) Text(content string) *CryoMessage {
	m.Elements = append(m.Elements, &TextElement{
		*lagrangeMessage.NewText(content),
	})
	return m // 返回当前消息对象，以便链式调用
}
func (m *CryoMessage) Texts(content ...string) *CryoMessage {
	if len(content) == 0 {
		return m // 如果没有内容，直接返回当前消息对象
	}
	for _, c := range content {
		m.Elements = append(m.Elements, &TextElement{
			*lagrangeMessage.NewText(c),
		})
	}
	return m
}
func (m *CryoMessage) At(uin uint32, display ...string) *CryoMessage {
	m.Elements = append(m.Elements, &AtElement{
		*lagrangeMessage.NewAt(uin, display...),
	})
	return m
}
func (m *CryoMessage) Face(faceId uint32) *CryoMessage {
	m.Elements = append(m.Elements, &FaceElement{
		*lagrangeMessage.NewFace(faceId),
	})
	return m
}
func (m *CryoMessage) Reply(msg CryoMessageEvent) *CryoMessage {
	replySeq, senderUin, time, elements := msg.replyDetail()
	m.Elements = append(m.Elements, &ReplyElement{
		lagrangeMessage.ReplyElement{
			ReplySeq:  replySeq,
			SenderUin: senderUin,
			Time:      time,
			Elements:  elements,
		},
	})
	return m
}

func (m *CryoMessage) Image(imgData []byte, summary ...string) *CryoMessage {
	m.Elements = append(m.Elements, &ImageElement{
		*lagrangeMessage.NewImage(imgData),
	})
	return m
}

func (m *CryoMessage) ImageIO(r io.ReadSeeker, summary ...string) *CryoMessage {
	m.Elements = append(m.Elements, &ImageElement{
		*lagrangeMessage.NewStreamImage(r, summary...),
	})
	return m
}

func (m *CryoMessage) ImageFile(filePath string, summary ...string) *CryoMessage {
	imgElement, err := lagrangeMessage.NewFileImage(filePath, summary...)
	if err != nil {
		Errorf("打开位于 %s 的图片时失败: %v", filePath, err)
	}
	m.Elements = append(m.Elements, &ImageElement{
		*imgElement,
	})
	return m
}

func (m *CryoMessage) Dice(value uint32) *CryoMessage {
	m.Elements = append(m.Elements, &FaceElement{
		*lagrangeMessage.NewDice(value),
	})
	return m
}
