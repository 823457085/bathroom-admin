package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	App      AppConfig      `yaml:"app"`
	JWT      JWTConfig      `yaml:"jwt"`
	WxPay    WxPayConfig    `yaml:"wxpay"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port    int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type AppConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type WxPayConfig struct {
	AppID     string `yaml:"app_id"`
	MchID     string `yaml:"mch_id"`
	APIKey    string `yaml:"api_key"`
	NotifyURL string `yaml:"notify_url"`
}

var GlobalConfig *Config

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	GlobalConfig = &cfg
	return &cfg, nil
}
