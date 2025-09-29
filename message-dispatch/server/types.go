package server

type MessageData struct {
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

type Message struct {
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageID   int64         `json:"message_id"`
	GroupID     int64         `json:"group_id"`
	UserID      int64         `json:"user_id"`
	Message     []MessageData `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Sender      Sender        `json:"sender"`
	Time        int64         `json:"time"`
	PostType    string        `json:"post_type"`
	// 其他字段可按需添加
}
