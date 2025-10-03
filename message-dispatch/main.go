package main

import (
	"fmt"
	"os"

	infra "github.com/li1553770945/onebot-agent-message-dispatch/infra/container"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	infra.InitGlobalContainer(env)
	App := infra.GetGlobalContainer()
	App.HttpServer.Start()
	fmt.Print(App.Config.Env)
}
