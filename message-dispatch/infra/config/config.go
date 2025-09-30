package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/li1553770945/onebot-agent-message-dispatch/infra/constant"
	"gopkg.in/yaml.v3"
)

type ReceiveRule struct {
	Name     string `yaml:"name"`      // 规则名称，便于识别
	FromType string `yaml:"from_type"` // 消息来源类型: private, group, all
	GroupId  string `yaml:"group_id"`  // 群号, 如果是私聊则不会处理，可以使用正则匹配，例如 "*"表示所有群
	UserId   string `yaml:"user_id"`   // 发送消息的用户号, 可以使用正则匹配，例如 "*" 表示所有用户
	ToAddr   string `yaml:"to_addr"`   // 目标服务
	IsEnd    bool   `yaml:"is_end"`    // 是否终止规则匹配
}
type ReceiverConfig struct {
	Rules []ReceiveRule `yaml:"rules"`
}

type SendRule struct {
	Name       string `yaml:"name"`        // 规则名称，便于识别
	OnebotType string `yaml:"onebot_type"` // onebot类型: lagrange, go-cqhttp, all
	SelfId     string `yaml:"self_id"`     // 机器人账号
	ToAddr     string `yaml:"to_addr"`     // 目标服务
	IsEnd      bool   `yaml:"is_end"`      // 是否终止规则匹配
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
	if env != constant.EnvProduction && env != constant.EnvDevelopment {
		panic(fmt.Sprintf("环境必须是%s或者%s之一", constant.EnvProduction, constant.EnvDevelopment))
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
