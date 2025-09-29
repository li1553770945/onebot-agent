package main

import (
	"fmt"
	infra "github.com/li1553770945/onebot-agent-message-dispatch/infra"
	"os"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	infra.InitGlobalContainer(env)
	App := infra.GetGlobalContainer()
	fmt.Print(App.Config.Env)
}
