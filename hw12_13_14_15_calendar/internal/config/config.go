package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Env     string `yaml:"env"`
	Logger  LoggerConf
	Server  ServerConf
	Storage StorageConf
}

type EnvConf struct {
	Level string
}

type LoggerConf struct {
	Level string
}

type ServerConf struct {
	Host        string
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type StorageConf struct {
	Type string
	SQL  SQLConf
}

type SQLConf struct {
	DSN string
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return &Config{}, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}
