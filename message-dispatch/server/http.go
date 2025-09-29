package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/li1553770945/onebot-agent-message-dispatch/infra/config"
)

type HttpServer struct {
	router *gin.Engine
	config *config.Config
}

func NewHttpServer(config *config.Config) *HttpServer {
	return &HttpServer{
		router: gin.Default(),
		config: config,
	}
}

func (s *HttpServer) Start() {
	r := s.router
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/", func(c *gin.Context) {
		s.HandleMessage(c)
	})
	err := r.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *HttpServer) HandleMessage(c *gin.Context) {
	bodyBytes, err := c.GetRawData()
	if err != nil {
		fmt.Println("读取 body 失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取 body 失败"})
		return
	}
	var msg Message
	if err := json.Unmarshal(bodyBytes, &msg); err != nil {
		fmt.Println("解析失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析失败"})
		return
	}
	s.HandleDispatchMessage(&msg)
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (s *HttpServer) HandleDispatchMessage(msg *Message) {
	if msg.PostType != "message" {
		return
	}
	matchCount := 0
	for _, rule := range s.config.Receiver.Rules {
		if s.IsMatchRule(msg, &rule) {
			fmt.Printf("匹配到规则: %+v\n", rule.ToService)
			matchCount++
			if rule.IsEnd {
				break
			}
		}
	}
	fmt.Println("匹配到规则数量:", matchCount)
}

func (s *HttpServer) IsMatchRule(msg *Message, rule *config.ReceiveRule) bool {
	// 匹配消息类型
	if rule.FromType != "all" && rule.FromType != msg.MessageType {
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
