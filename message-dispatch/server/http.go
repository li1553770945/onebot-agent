package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/li1553770945/onebot-agent-message-dispatch/infra/config"
	"github.com/li1553770945/onebot-agent-message-dispatch/server/sender"
)

type HttpServer struct {
	router         *gin.Engine
	config         *config.Config
	LagrangeSender *sender.LagrangeSender
}

func NewHttpServer(config *config.Config, lagrangeSender *sender.LagrangeSender) *HttpServer {
	return &HttpServer{
		router:         gin.Default(),
		config:         config,
		LagrangeSender: lagrangeSender,
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
		s.HandleReceiveMessage(c)
	})
	r.POST("/send", func(c *gin.Context) {
		s.HandleSendMessage(c)
	})
	err := r.Run(":15001")
	if err != nil {
		fmt.Println(err)
		return
	}
}
