package infra

import (
	"sync"
)

type Container struct {
	Config *Config
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

func NewContainer(config *Config) *Container {
	return &Container{
		Config: config,
	}

}

func GetContainer(env string) *Container {
	config := GetConfig(env)
	app := NewContainer(config)
	return app
}
