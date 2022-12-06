package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Log              string     `yaml:"log"`
	LogLevel         string     `yaml:"log_level"`
	ListenPort       int        `yaml:"listen_port"`
	WeChatApiAddress string     `yaml:"wechat_api_address"`
	MessageNotifyURL []string   `yaml:"message_notify_url"`
	QunManager       QunManager `yaml:"qun_manager"`
	Module           Module     `yaml:"module"`
}

type Module struct {
	ChatGPT ChatGPTConf `yaml:"chatgpt"`
}

type ChatGPTConf struct {
	AiSession string `yaml:"ai_session"`
}

type QunManager struct {
	EnableWord  []string `yaml:"enable_word"`
	DisbaleWord []string `yaml:"disbale_word"`
}

func Load(path string) (*Config, error) {
	value, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rs := &Config{}
	if err := yaml.Unmarshal(value, rs); err != nil {
		return nil, err
	}
	return rs, nil
}
