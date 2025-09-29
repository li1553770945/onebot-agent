package config

import (
	"fmt"
	"github.com/li1553770945/onebot-agent-message-dispatch/infra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type ReceiveRule struct {
	FromId    string `yaml:"from_id"`    // 可以是群号或者用户号
	ToService string `yaml:"to_service"` // 目标服务
}
type ReceiverConfig struct {
	Rules []ReceiveRule `yaml:"rules"`
}

type SendRule struct {
	SelfId    string `yaml:"self_id"`    // 机器人账号
	ToService string `yaml:"to_service"` // 目标服务
}
type SenderConfig struct {
	Rules []SendRule `yaml:"rules"`
}

type Config struct {
	Receiver ReceiverConfig `yaml:"receiver"`
	Sender   SenderConfig   `yaml:"sender"`
	Env      string         `yaml:"env"`
}

func GetConfig(env string) *Config {
	if env != infra.EnvProduction && env != infra.EnvDevelopment {
		panic(fmt.Sprintf("环境必须是%s或者%s之一", infra.EnvProduction, infra.EnvDevelopment))
	}
	conf := &Config{}
	path := filepath.Join("config", fmt.Sprintf("%s.yml", env))
	fmt.Println("path:", path)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		panic("配置文件不存在")
	}
	if err != nil {
		panic(err)
	}
	err = yaml.NewDecoder(f).Decode(conf)
	conf.Env = env
	if err != nil {
		panic(err)
	}

	return conf
}
