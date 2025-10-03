package types

type SendMessageParams struct {
	DetailType string        `json:"detail_type"` // 消息子类型: private, group
	UserID     string        `json:"user_id"`     // 私聊发送消息的用户号
	GroupID    string        `json:"group_id"`    // 群聊发送消息的群号
	Message    []interface{} `json:"message"`
	SelfID     string        `json:"self_id"` // 机器人账号
}
type SendMessage struct {
	Action string            `json:"action"`
	Params SendMessageParams `json:"params"`
}
