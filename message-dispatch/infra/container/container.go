package container

import (
	"sync"

	"github.com/li1553770945/onebot-agent-message-dispatch/infra/config"
	"github.com/li1553770945/onebot-agent-message-dispatch/server"
	"github.com/li1553770945/onebot-agent-message-dispatch/server/sender"
)

type Container struct {
	Config     *config.Config
	HttpServer *server.HttpServer
}

var APP *Container
var once sync.Once

func GetGlobalContainer() *Container {
	if APP == nil {
		panic("APP在使用前未初始化")
	}
	return APP
}

func InitGlobalContainer(env string) {
	once.Do(func() {
		APP = GetContainer(env)
	})
}

func NewContainer(config *config.Config,
	httpServer *server.HttpServer,
) *Container {
	return &Container{
		Config:     config,
		HttpServer: httpServer,
	}

}

func GetContainer(env string) *Container {
	config := config.GetConfig(env)

	lagrangeSender := sender.NewLagrangeSender()
	httpServer := server.NewHttpServer(config, lagrangeSender)
	app := NewContainer(config, httpServer)
	return app
}
