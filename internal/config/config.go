package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Log              string   `yaml:"log"`
	LogLevel         string   `yaml:"log_level"`
	ListenPort       int      `yaml:"listen_port"`
	MessageNotifyURL []string `yaml:"message_notify_url"`
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
