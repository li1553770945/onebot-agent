package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/li1553770945/onebot-agent-message-dispatch/infra/config"
	"github.com/li1553770945/onebot-agent-message-dispatch/infra/constant"
	"github.com/li1553770945/onebot-agent-message-dispatch/server/types"
)

func (s *HttpServer) HandleSendMessage(c *gin.Context) {
	bodyBytes, err := c.GetRawData()
	if err != nil {
		fmt.Println("读取 body 失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取 body 失败"})
		return
	}
	var msg types.SendMessage
	if err := json.Unmarshal(bodyBytes, &msg); err != nil {
		fmt.Println("解析失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析失败"})
		return
	}
	s.HandleSendDispatchMessage(&msg)
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (s *HttpServer) HandleSendDispatchMessage(msg *types.SendMessage) {
	if msg.Action != "send_message" {
		return
	}
	matchCount := 0
	for _, rule := range s.config.Sender.Rules {
		if s.IsSendMatchRule(msg, &rule) {
			fmt.Printf("匹配到规则: %+v\n", rule.Name)
			s.SendMessage(&rule, msg)
			matchCount++
			if rule.IsEnd {
				break
			}
		}
	}
	if matchCount == 0 {
		fmt.Println("未匹配到任何规则，消息不会被发送")
	} else {
		fmt.Println("匹配到规则数量:", matchCount)
	}

}

func (s *HttpServer) SendMessage(rule *config.SendRule, msg *types.SendMessage) {
	if rule.OnebotType == constant.LagrangeOnebotType {
		s.LagrangeSender.SendMessage(rule.ToAddr, msg)
	} else {
		fmt.Printf("不支持的发送端类型:%s，该消息不会被发送\n", rule.OnebotType)
	}
}
func (s *HttpServer) IsSendMatchRule(msg *types.SendMessage, rule *config.SendRule) bool {
	// 匹配消息类型
	if rule.SelfId == msg.Params.SelfID {
		return true
	}
	return false
}
