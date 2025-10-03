package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/li1553770945/onebot-agent-message-dispatch/server/types"
)

type PrivateSendReq struct {
	MessageType string      `json:"message_type"` // 消息类型: private, group
	UserID      string      `json:"group_id"`     // 群 Uin
	Message     interface{} `json:"message"`      // 消息内容（可以是消息段对象、消息段数组或字符串）
}
type GroupSendReq struct {
	MessageType string      `json:"message_type"` // 消息类型: private, group
	GroupID     string      `json:"group_id"`     // 群 Uin
	Message     interface{} `json:"message"`      // 消息内容（可以是消息段对象、消息段数组或字符串）
}

type LagrangeSender struct {
}

func NewLagrangeSender() *LagrangeSender {
	return &LagrangeSender{}
}

func (l *LagrangeSender) SendMessage(to_service string, msg *types.SendMessage) {
	if msg.Params.DetailType == "group" {
		l.SendGroupMessage(to_service, msg)
	}
	if msg.Params.DetailType == "private" {
		l.SendPrivateMessage(to_service, msg)
	}
}

func (l *LagrangeSender) SendGroupMessage(to_service string, msg *types.SendMessage) {
	addr := fmt.Sprint(to_service, "/send_msg")
	req := &GroupSendReq{
		MessageType: "group",
		GroupID:     msg.Params.GroupID,
		Message:     msg.Params.Message,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("消息序列化失败:", err)
		return
	}
	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("发送消息失败:", err)
		return
	}
}
func (l *LagrangeSender) SendPrivateMessage(to_service string, msg *types.SendMessage) {
	addr := fmt.Sprint(to_service, "/send_msg")
	req := &PrivateSendReq{
		MessageType: "private",
		UserID:      msg.Params.UserID,
		Message:     msg.Params.Message,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("消息序列化失败:", err)
		return
	}
	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("发送消息失败:", err)
		return
	}
}
