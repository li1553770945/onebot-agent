package server

import (
	"fmt"
	"net/http"

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
		bodyBytes, err := c.GetRawData()
		if err != nil {
			fmt.Println("读取 body 失败:", err)
		} else {
			fmt.Println(string(bodyBytes))
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	err := r.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
