package types

import (
	"bytes"
	"encoding/json"
)

// StringOrNumber 用于兼容 OneBot 上游实现里既可能是字符串也可能是数字的 ID 字段。
// 反序列化时：如果是字符串，直接取内容；如果是数字，转为其十进制字符串表示。
type StringOrNumber string

func (s *StringOrNumber) UnmarshalJSON(data []byte) error {
	// null 视为空字符串
	if bytes.Equal(data, []byte("null")) {
		*s = ""
		return nil
	}
	// 如果以引号开头，当作正常字符串处理
	if len(data) > 0 && data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		*s = StringOrNumber(str)
		return nil
	}
	// 尝试按数字解析
	var num json.Number
	if err := json.Unmarshal(data, &num); err != nil {
		return err
	}
	*s = StringOrNumber(num.String())
	return nil
}

// ToString 便于外部转换为普通 string
func (s StringOrNumber) ToString() string { return string(s) }

type ReceiveMessageData struct {
	Type string `json:"type"`
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

type Sender struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	// 其他字段可按需添加
}

type ReveiceMessage struct {
	MessageType string               `json:"message_type"`
	SubType     string               `json:"sub_type"`
	MessageID   int64                `json:"message_id"`
	GroupID     int64                `json:"group_id"`
	UserID      int64                `json:"user_id"`
	Message     []ReceiveMessageData `json:"message"`
	RawMessage  string               `json:"raw_message"`
	Sender      Sender               `json:"sender"`
	Time        int64                `json:"time"`
	PostType    string               `json:"post_type"`
	SelfID      StringOrNumber       `json:"self_id"`
	// 其他字段可按需添加
}
