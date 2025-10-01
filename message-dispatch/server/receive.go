package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/li1553770945/onebot-agent-message-dispatch/infra/config"
	"github.com/li1553770945/onebot-agent-message-dispatch/server/types"
)

func (s *HttpServer) HandleReceiveMessage(c *gin.Context) {
	bodyBytes, err := c.GetRawData()

	if err != nil {
		fmt.Println("读取 body 失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取 body 失败"})
		return
	}
	var msg types.ReveiceMessage
	if err := json.Unmarshal(bodyBytes, &msg); err != nil {
		fmt.Println("解析失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析失败"})
		return
	}
	if msg.PostType != "meta_event" {
		fmt.Printf("收到消息：%s\n", string(bodyBytes))
		s.HandleReceiveDispatchMessage(&msg, bodyBytes)
	}
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (s *HttpServer) HandleReceiveDispatchMessage(msg *types.ReveiceMessage, bodyBytes []byte) {
	matchCount := 0
	for _, rule := range s.config.Receiver.Rules {
		if s.IsReceiveMatchRule(msg, &rule) {
			fmt.Printf("匹配到规则: %+v\n", rule.Name)
			s.SendToService(&rule, bodyBytes)
			matchCount++
			if rule.IsEnd {
				break
			}
		}
	}
	fmt.Println("匹配到规则数量:", matchCount)
}
func (s *HttpServer) SendToService(rule *config.ReceiveRule, bodyBytes []byte) {
	resp, err := http.Post(rule.ToAddr, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("Receive发送消息到服务失败: %v,HTTP状态码: %d", err, resp.StatusCode)
		return
	}

}
func (s *HttpServer) IsReceiveMatchRule(msg *types.ReveiceMessage, rule *config.ReceiveRule) bool {
	// 匹配消息类型
	if rule.FromType != "all" && rule.FromType != msg.MessageType && rule.FromType != msg.RequestType {
		return false
	}
	// 匹配 group_id（群聊时才判断）
	if msg.MessageType == "group" {
		matched, err := regexp.MatchString(rule.GroupId, fmt.Sprint(msg.GroupID))
		if err != nil || !matched {
			return false
		}
	}
	// 匹配 user_id
	matched, err := regexp.MatchString(rule.UserId, fmt.Sprint(msg.UserID))
	if err != nil || !matched {
		return false
	}
	return true
}
